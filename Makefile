# Makefile
.PHONY: build

BINARY_NAME=pubsbot

build:
	go mod tidy && \
	go build -ldflags="-w -s" -o ./bin/${BINARY_NAME}${SUFFIX} ./cmd/pubsbot

run:
	make build
	bin/pubsbot

