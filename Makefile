SRC_FILES = $(shell find . -name '*.go' -not -path './vendor/*')

IMPORT_PATH = go.alexhamlin.co/slackbridge
VERSION = $(shell git describe | sed 's/^v//')

.PHONY: release
release: slackbridge-linux-amd64

slackbridge-linux-amd64: $(SRC_FILES)
	GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build \
		-mod vendor \
		-o slackbridge-linux-amd64 \
		-ldflags "-s -w -X $(IMPORT_PATH)/cmd.Version=$(VERSION)" \
		.

.PHONY: clean
clean:
	rm -f slackbridge-linux-amd64
