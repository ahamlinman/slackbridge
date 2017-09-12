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
	"os"
	"os/exec"
	"sync"
)

// Process is the type for a child process managed by package childproc.
type Process struct {
	process      *os.Process
	shutdownOnce sync.Once
	err          error
}

// Spawn starts a child process from the given command line (name + arguments)
// whose standard streams are connected to the provided io.ReadCloser and
// io.WriteCloser. After the child process terminates (either on its own or
// through explicit action from the caller), or if an error is encountered
// while spawning the process, the provided ReadCloser and WriteCloser will be
// closed.
func Spawn(cmdline []string, stdin io.ReadCloser, stdouterr io.WriteCloser) (*Process, error) {
	// Look up based on $PATH, just like package exec
	path := exec.LookPath(cmdline[0])

	attr := &os.ProcAttr{
		Files: []*File{}, // TODO implement this for real
	}

	proc, err := os.StartProcess(path, cmdline[1:], attr)
	if err != nil {
		return nil, fmt.Errorf("childproc failed to Spawn: %v", err)
	}

	return &Process{
		process: proc,
	}
}

// Close terminates this child process if it has not yet terminated on its own.
// It also returns any error encountered as a result of running and/or
// terminating the process.
func (p *Process) Close() error {
	p.shutdown(true)
	return p.err
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
	})
}
