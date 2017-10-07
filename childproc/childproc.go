/*

Package childproc implements spawning of child processes whose standard streams
are connected to an io.ReadCloser and io.WriteCloser.

Specifically, childproc differs from the standard os/exec package in its
handling of process termination. The Wait method of exec.Cmd waits for any
io.Reader connected to the stdin of the process to terminate. However, we may
wish to close an unbounded Reader only after the connected child process
terminates. This creates a circular dependency.

childproc resolves this by owning the closure of the provided Reader and
Writer, closing both after the process terminates. Notably, it assumes that
calling Close on the provided io.ReadCloser will interrupt an active concurrent
Read, causing it to return EOF. While there is some indication (e.g.
https://stackoverflow.com/a/26441866) that other io.ReadCloser implementations
do this, childproc's behavior is not necessarily expected to be useful for
consumers other than slackbridge and readers other than that provided by
slackio. This assumption regarding Close behavior is certainly not guaranteed
for arbitrary readers (e.g. OS pipes).

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
	errCh        chan error
	err          error
}

// numWaitErrors is the number of errors that Spawn might send into the errCh
// of a Process.
var numWaitErrors = 6

// Spawn starts a child process from the given command line (name + arguments)
// whose standard streams are connected to the provided io.ReadCloser and
// io.WriteCloser. If the child process is started successfully (err == nil),
// the provided ReadCloser and WriteCloser will be closed after it terminates.
// Otherwise they will be left open.
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

	attrs := &os.ProcAttr{
		Files: []*os.File{childStdinOut, childStdoutIn, childStdoutIn},
	}
	process, err := os.StartProcess(path, cmdline, attrs)
	if err != nil {
		return nil, fmt.Errorf("childproc start failed: %v", err)
	}

	// Note that from here on out, we no longer return with err != nil. We need
	// to fulfill our documented contract of closing inputReader and outputWriter
	// when the child terminates. So errors go into this channel to be returned
	// by Wait.
	//
	// NOTE: Don't forget to update numWaitErrors if a new send into this channel
	// is added!
	errCh := make(chan error, numWaitErrors)

	// See comments in Wait for why we don't also close childStdinOut
	errCh <- childStdoutIn.Close() // Send 1 into errCh

	// Spawn copy goroutines for the provided reader and writer
	// First, from the reader to the child stdin
	go func() {
		_, copyErr := io.Copy(childStdinIn, inputReader)
		errCh <- copyErr // 2

		// inputReader closed by Wait after child terminates
		errCh <- childStdinIn.Close() // 3
	}()

	// Second, from the child stdout to the writer
	go func() {
		_, copyErr := io.Copy(outputWriter, childStdoutOut)
		errCh <- copyErr // 4

		errCh <- childStdoutOut.Close() // 5
		errCh <- outputWriter.Close()   // 6
	}()

	p := &Process{
		process:       process,
		readerCloser:  inputReader,
		childStdinOut: childStdinOut,
		errCh:         errCh,
	}

	// Ensure we clean up regardless of whether consumers call Wait
	go p.Wait()

	return p, nil
}

// Wait waits for the process created by Spawn to terminate, and returns any
// errors encountered while waiting on the process or copying to/from its
// standard streams.
//
// It is safe to call Wait more than once and/or concurrently. All calls will
// return the same error value.
func (p *Process) Wait() error {
	p.shutdownOnce.Do(func() {
		var errs *multierror.Error

		_, err := p.process.Wait()
		errs = multierror.Append(errs, err)

		// A goroutine feeds the Reader's output to the child's stdin through a
		// pipe. Because that goroutine could block on writing to the pipe, we
		// drain the output side of that pipe ourselves. This drain operation will
		// finish when the goroutine closes the input side of the pipe.
		//
		// This setup is not ideal, since we need to keep a reference to the output
		// side of the stdin pipe. If we handled EPIPE in a cross-platform way like
		// package exec, slackbridge could spawn twice as many processes without
		// hitting limits on open files. However, this would warrant cross-platform
		// testing that I can't easily do.
		errs = multierror.Append(errs, p.readerCloser.Close())
		_, err = io.Copy(ioutil.Discard, p.childStdinOut)
		errs = multierror.Append(errs, err, p.childStdinOut.Close())
		// We do *not* keep a reference to the input side of the stdout pipe, so
		// termination of the child process will EOF the output side and let that
		// goroutine stop.

		for i := 0; i < numWaitErrors; i++ {
			errs = multierror.Append(errs, <-p.errCh)
		}

		p.err = errs.ErrorOrNil()
	})

	return p.err
}
