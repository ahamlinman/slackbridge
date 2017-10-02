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
	"sync"

	multierror "github.com/hashicorp/go-multierror"
)

// Process is the type for a child process managed by package childproc.
type Process struct {
	process *os.Process

	readerCloser  io.Closer
	childStdinOut *os.File

	shutdownOnce sync.Once
	stdinErrCh   chan error
	stdoutErrCh  chan error
	err          error
}

// Spawn starts a child process from the given command line (name + arguments)
// whose standard streams are connected to the provided io.ReadCloser and
// io.WriteCloser. If the child process is started successfully, the provided
// ReadCloser and WriteCloser will be closed after it terminates.
func Spawn(cmdline []string, inputReader io.ReadCloser, outputWriter io.WriteCloser) (proc *Process, err error) {
	// Look up based on $PATH, just like package exec
	path, err := exec.LookPath(cmdline[0])
	if err != nil {
		return nil, fmt.Errorf("childproc lookup failed: %v", err)
	}

	// Create OS pipes for standard streams
	// First, for the child's stdin
	childStdinOut, childStdinIn, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("childproc pipe creation failed: %v", err)
	}
	defer func() {
		if err != nil {
			childStdinOut.Close()
			childStdinIn.Close()
		}
	}()

	// Second, for the child's stdout
	childStdoutOut, childStdoutIn, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("childproc pipe creation failed: %v", err)
	}
	defer func() {
		if err != nil {
			childStdoutOut.Close()
			childStdoutIn.Close()
		}
	}()

	// Start the real child process
	attrs := &os.ProcAttr{
		Files: []*os.File{childStdinOut, childStdoutIn, childStdoutIn},
	}
	process, err := os.StartProcess(path, cmdline[1:], attrs)
	if err != nil {
		return nil, fmt.Errorf("childproc start failed: %v", err)
	}

	childStdoutIn.Close()

	stdinErrCh, stdoutErrCh := make(chan error), make(chan error)

	// Spawn copy goroutines for the provided reader and writer
	// Very, very manual. Much less fancy than package exec.
	// First, from the reader to the child stdin
	go func() {
		_, copyErr := io.Copy(childStdinIn, inputReader)
		childStdinIn.Close()
		stdinErrCh <- copyErr
	}()

	// Second, from the child stdout to the writer
	go func() {
		_, copyErr := io.Copy(outputWriter, childStdoutOut)
		childStdoutOut.Close()
		outputWriter.Close()
		stdoutErrCh <- copyErr
	}()

	p := &Process{
		process:       process,
		readerCloser:  inputReader,
		childStdinOut: childStdinOut,
		stdinErrCh:    stdinErrCh,
		stdoutErrCh:   stdoutErrCh,
	}

	go p.Wait()

	return p, nil
}

func (p *Process) Wait() error {
	p.shutdownOnce.Do(func() {
		var errs *multierror.Error

		// Wait on the process and get any errors from it
		_, err := p.process.Wait()
		errs = multierror.Append(errs, err)

		// A goroutine feeds the Reader's output to the child's stdin through a
		// pipe. Because that goroutine could block on writing to the pipe, we
		// drain the output side of in the pipe ourselves. This drain operation
		// will finish when the goroutine closes the input side of the pipe,
		// allowing us to collect any error emitted by the goroutine.
		//
		// This setup is not ideal, since we need to keep a reference to the output
		// side of the stdin pipe. If we handled EPIPE in a cross-platform way like
		// package exec, slackbridge could spawn twice as many processes without
		// hitting limits on open files.
		errs = multierror.Append(errs, p.readerCloser.Close())
		_, err = io.Copy(ioutil.Discard, p.childStdinOut)
		errs = multierror.Append(errs, err, <-p.stdinErrCh, p.childStdinOut.Close())

		// We do *not* keep a reference to the input side of the stdout pipe, so
		// termination of the child process will EOF the output side and let that
		// goroutine stop. This ensures that we can safely shut down the writer.
		errs = multierror.Append(errs, <-p.stdoutErrCh)

		p.err = errs.ErrorOrNil()
	})

	return p.err
}
