package slackio

// Message is the type for messages received from and sent to a single Slack
// channel.
type Message struct {
	ID        int
	ChannelID string
	Text      string
}
