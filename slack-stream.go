package main

import (
	"github.com/nlopes/slack"
)

// slackLineStreamer streams message lines (as strings) from a given Slack
// channel in real-time, and provides an API for replies.
type slackLineStreamer struct {
	rtm          *slack.RTM
	slackChannel string
	lineStream   chan string
	close        chan bool
}

func newSlackLineStreamer(apiToken, channel string) *slackLineStreamer {
	api := slack.New(apiToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	sls := &slackLineStreamer{
		rtm:          rtm,
		slackChannel: channel,
		lineStream:   make(chan string),
		close:        make(chan bool),
	}
	go sls.readLineStream()

	return sls
}

func (s *slackLineStreamer) readLineStream() {
	for {
		select {
		case evt := <-s.rtm.IncomingEvents:
			if data, ok := evt.Data.(*slack.MessageEvent); ok {
				if data.Type == "message" && data.Channel == s.slackChannel && data.Text != "" {
					s.lineStream <- data.Text
				}
			}

		case <-s.close:
			return
		}
	}
}

func (s *slackLineStreamer) ReceiveChan() <-chan string {
	return s.lineStream
}

func (s *slackLineStreamer) Send(text string) {
	msg := s.rtm.NewOutgoingMessage(text, s.slackChannel)
	s.rtm.SendMessage(msg)
}

func (s *slackLineStreamer) Close() error {
	s.close <- true
	return s.rtm.Disconnect()
}
