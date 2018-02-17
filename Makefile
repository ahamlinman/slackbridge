SRC_FILES = $(shell find . -name '*.go' -not -path './vendor/*')

IMPORT_PATH = github.com/ahamlinman/slackbridge
VERSION = $(shell git describe | sed 's/^v//')

.PHONY: release
release: slackbridge-linux-amd64

slackbridge-linux-amd64: $(SRC_FILES) vendor
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build \
		-o slackbridge-linux-amd64 \
		-ldflags "-s -w -X $(IMPORT_PATH)/cmd.Version=$(VERSION)" \
		.

vendor Gopkg.lock: $(SRC_FILES) Gopkg.toml
	dep ensure -v
	touch vendor

.PHONY: test
test: vendor
	go test ./slackio -cover -timeout 10s -race

.PHONY: clean
clean:
	rm -f slackbridge-linux-amd64
