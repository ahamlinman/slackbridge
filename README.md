# slackbridge

**slackbridge** is designed to connect Slack channels to executable programs on
a server, by translating real-time messages to and from standard I/O channels.

At this stage, it's a very very rough prototype. However, it does generally do
most of its job. You can configure it to run a program (e.g. the "ed" text
editor) and pretty much use that program via Slack.

### TODOs

* Clean up execLineStreamer
  - Automatically re-spawn the program on exit
  - Better error handling
* Fix processing of last Slack message on startup (seems that Slack or the
  library I'm using tries to treat that message as new, which I don't want)
