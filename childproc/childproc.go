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
)

// Process is the type for a child process managed by package childproc.
type Process struct {
	process *os.Process

	readerCloser io.Closer
	writerCloser io.Closer

	childInRead   *os.File // reference kept to drain pipe after child exit
	childInWrite  *os.File
	childOutRead  *os.File
	childOutWrite *os.File

	shutdownOnce sync.Once
	err          error
}

// Spawn starts a child process from the given command line (name + arguments)
// whose standard streams are connected to the provided io.ReadCloser and
// io.WriteCloser. After the child process terminates (either on its own or
// through explicit action from the caller) the provided ReadCloser and
// WriteCloser will be closed.
func Spawn(cmdline []string, stdin io.ReadCloser, stdouterr io.WriteCloser) (*Process, error) {
	p := &Process{
		readerCloser: stdin,
		writerCloser: stdouterr,
	}

	// Look up based on $PATH, just like package exec
	path, err := exec.LookPath(cmdline[0])
	if err != nil {
		return nil, fmt.Errorf("childproc lookup failed: %v", err)
	}

	if err := p.populateIOStreams(); err != nil {
		return nil, fmt.Errorf("chlidproc setup failed: %v", err)
	}

	attr := &os.ProcAttr{
		Files: []*os.File{p.childInWrite, p.childOutRead, p.childOutRead},
	}

	proc, err := os.StartProcess(path, cmdline[1:], attr)
	if err != nil {
		// TODO close stuff
		return nil, fmt.Errorf("childproc start failed: %v", err)
	}

	p.process = proc
	return p, nil
}

func (p *Process) populateIOStreams() error {
	// TODO stop leaking fds, but make it clean

	childInRead, childInWrite, err := os.Pipe()
	if err != nil {
		return err
	}

	childOutRead, childOutWrite, err := os.Pipe()
	if err != nil {
		return err
	}

	p.childInRead, p.childInWrite = childInRead, childInWrite
	p.childOutRead, p.childOutWrite = childOutRead, childOutWrite
	return nil
}

// Close terminates this child process if it has not yet terminated on its own.
// It also returns any error encountered as a result of running and/or
// terminating the process.
func (p *Process) Close() error {
	p.shutdown(true)
	return p.err
}

// closeAllIO shuts down all I/O streams referenced by this Process, including
// the Reader and Writer provided when spawning as well as all internal pipes.
func (p *Process) closeAllIO() error {
	var err error

	// 1. Close the input Reader, so it will return EOF and copies out of it will
	// terminate.
	if p.readerCloser != nil {
		err = p.readerCloser.Close()
	}

	// 2. Close our connection to the child's stdin, draining the output to
	// prevent EPIPEs (kind of dumb?)
	if p.childInRead != nil {
		done := make(chan struct{})
		go func() {
			// TODO assign this to err?
			io.Copy(ioutil.Discard, p.childInRead)
			close(done)
		}()

		err = p.childInWrite.Close()
		<-done // TODO also wait for goroutine?
		err = p.childInRead.Close()
		p.childInWrite, p.childInRead = nil, nil
	}

	if p.childOutWrite != nil {
		err = p.childOutWrite.Close()
		// TODO wait for goroutine to terminate
		err = p.writerCloser.Close()
	}

	return err
}

// shutdown rolls up logic that should only run once when a child process
// terminates. kill determines whether the process will be explicitly killed if
// it has not terminated on its own already.
func (p *Process) shutdown(kill bool) {
	p.shutdownOnce.Do(func() {
		if kill {
			// TODO kill
			// this goes into the sync.Once instead of Close, because if the process
			// has already terminated on its own we don't want to try killing it
		}

		// TODO implement for real

		// 1. the process terminates (above)

		// 2. close io.Reader, stops io.Copy
		// 2a. drain child stdin? to be fair we will get an EPIPE if we don't
		// (meaning we have to reimplement per-OS os/exec logic)
		// (but yes we also need to retain a reference to the fd in this case)
		// 3. close os.Pipe write side (child stdin)

		// 4. close os.Pipe read side (child stdout/err), stops io.Copy with nil err
		// 5. close io.Writer
		// (behavior of pipe read close: https://play.golang.org/p/tq8QVRLKug)
	})
}
