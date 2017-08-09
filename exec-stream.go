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
	// Yes, this is probably super inefficient. And it ignores errors.
	s.stdin.Write([]byte(text))
	s.stdin.Write([]byte{'\n'})
}

func (s *execLineStreamer) Close() error {
	s.killCmd()
	s.cmd.Wait() // yes, ignoring error again
	return nil   // this is, like, really bad
}
