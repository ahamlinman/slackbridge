/*

Package childproc implements spawning of child processes whose standard streams
are connected to an io.ReadCloser and io.WriteCloser.

Specifically, childproc differs from the standard os/exec package in its
handling of process termination. The Wait method of exec.Cmd waits for any
io.Reader connected to the stdin of the process to terminate. However, when the
output of the Reader is unbounded, any request to terminate it may depend on
termination of the child process. This creates a circular dependency.

childproc resolves this dependency by assuming that the provided Reader and
Writer are unbounded and owning their closure after termination of the child
process connected to them. This behavior is not necessarily expected to be
useful for consumers other than slackbridge.

*/
package childproc

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// Process is the type for a child process managed by package childproc.
type Process struct {
	process *os.Process

	reader io.Closer
	writer io.Closer

	childStdinOut  *os.File // yes, we reference this, see below
	childStdoutOut *os.File

	stdinErr  chan error
	stdoutErr chan error
}

// Spawn starts a child process from the given command line (name + arguments)
// whose standard streams are connected to the provided io.ReadCloser and
// io.WriteCloser. If the child process is started successfully, the provided
// ReadCloser and WriteCloser will be closed after it terminates.
func Spawn(cmdline []string, stdin io.ReadCloser, stdouterr io.WriteCloser) (proc *Process, err error) {
	// Look up based on $PATH, just like package exec
	path, err := exec.LookPath(cmdline[0])
	if err != nil {
		err = fmt.Errorf("childproc lookup failed: %v", err)
		return
	}

	proc = &Process{
		stdinErr:  make(chan error),
		stdoutErr: make(chan error),
	}

	// Create OS pipes for standard streams
	// First, for the child's stdin
	var childStdinIn *os.File
	proc.childStdinOut, childStdinIn, err = os.Pipe()
	if err != nil {
		err = fmt.Errorf("childproc pipe creation failed: %v", err)
		return
	}
	defer func() {
		if err != nil {
			// Ignoring further errors...
			childStdinIn.Close()
			proc.childStdinOut.Close()
		}
	}()

	// Second, for the child's stdout
	var childStdoutIn *os.File
	proc.childStdoutOut, childStdoutIn, err = os.Pipe()
	if err != nil {
		err = fmt.Errorf("childproc pipe creation failed: %v", err)
		return
	}
	defer func() {
		if err != nil {
			// Again, ignoring further errors...
			childStdoutIn.Close()
			proc.childStdoutOut.Close()
		}
	}()

	// Start the real child process
	attrs := &os.ProcAttr{
		Files: []*os.File{proc.childStdinOut, childStdoutIn, childStdoutIn},
	}
	proc.process, err = os.StartProcess(path, cmdline[1:], attrs)
	if err != nil {
		err = fmt.Errorf("childproc start failed: %v", err)
		return
	}

	childStdoutIn.Close()

	// Spawn copy goroutines for the provided reader and writer
	// Very, very manual. Much less fancy than package exec.
	// First, from the reader to the child stdin
	go func() {
		_, copyErr := io.Copy(childStdinIn, stdin)
		childStdinIn.Close()
		proc.stdinErr <- copyErr
	}()

	// Second, from the child stdout to the writer
	go func() {
		_, copyErr := io.Copy(stdouterr, proc.childStdoutOut)
		stdouterr.Close()
		proc.stdoutErr <- copyErr
	}()

	return
}

func (p *Process) Wait() error {
	// Wait on the process and get any errors from it
	_, err := p.process.Wait()

	// Shut down the reader, to stop its copy goroutine with EOF
	err = p.reader.Close()
	// Drain the rest of the child's stdin pipe to ioutil.Discard
	_, err = io.Copy(ioutil.Discard, p.childStdinOut)
	// Obtain the error result from the copy goroutine
	err = <-p.stdinErr
	// Shut down our reference to the stdin read fd
	err = p.childStdinOut.Close()

	// Obtain the error result from the copy goroutine (child closure shuts down)
	err = <-p.stdoutErr
	// Shut down our read end of the child pipe, since it gets no more data
	err = p.childStdoutOut.Close()

	// Somehow return whatever error makes the most sense
	return err
}
