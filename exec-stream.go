package main

import (
	"bufio"
	"context"
	"io"
	"os/exec"
)

// execLineStreamer streams lines from a running program in real time.
type execLineStreamer struct {
	cmd        *exec.Cmd
	stdin      io.Writer
	stdout     io.Reader
	killCmd    context.CancelFunc
	lineStream chan string
}

func newExecLineStreamer(execLine []string) *execLineStreamer {
	ctx, killCmd := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, execLine[0], execLine[1:]...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// At some point, we need to call Wait in order to clean up the program's
	// resources after it exits. This is really easy to do in a goroutine though.

	// My thought right now is to replace the stdio pipe setup here with my own
	// custom pipes (disconnected from exec.Cmd), and have a goroutine that
	// re-spawns the program as necessary, probably with new pipes. Still need to
	// figure out how I can do that in a way that keeps the lineStream working
	// properly.

	// I probably need to separate out the lineStream from all of the command
	// execution stuff. Just rebuild some kind of command struct as necessary and
	// keep passing a single lineStream to it.

	els := &execLineStreamer{
		cmd:        cmd,
		stdin:      stdin,
		stdout:     stdout,
		killCmd:    killCmd,
		lineStream: make(chan string),
	}
	go els.readLineStream()

	return els
}

func (s *execLineStreamer) readLineStream() {
	scanner := bufio.NewScanner(s.stdout)
	for scanner.Scan() {
		s.lineStream <- scanner.Text()
	}
}

func (s *execLineStreamer) ReceiveChan() <-chan string {
	return s.lineStream
}

func (s *execLineStreamer) Send(text string) {
	// This ignores errors. Still need to figure out how to deal with that.
	s.stdin.Write(append([]byte(text), byte('\n')))
}

func (s *execLineStreamer) Close() error {
	s.killCmd()
	s.cmd.Wait() // yes, ignoring error again
	return nil   // this is, like, really bad
}
