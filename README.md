# slackbridge

**slackbridge** is designed to connect Slack channels to executable programs on
a server, by translating real-time messages to and from standard I/O channels.

At this stage, it's a very very rough prototype. However, it does generally do
most of its job. You can configure it to run a program (e.g. the "ed" text
editor) and pretty much use that program via Slack.

### TODOs

* Clean up execLineStreamer
  - Better error handling
  - Handle the case of the program exiting
* General architectural and code style points: I wrote this in a matter of a
  few hours over an evening (including time spent setting up Slack and learning
  its API), and frankly my mind was a bit dead so I totally stumbled through
  the whole thing. I'm not sure how great this "LineStreamer" concept really
  is. Like, maybe I should just make everything a ReadWriter?
