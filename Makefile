PROJECT=pushd-client-go

BUILD_PATH := $(shell pwd)/.gobuild

C0_PATH := $(BUILD_PATH)/src/github.com/catalyst-zero

BIN := $(PROJECT)

.PHONY=clean run-test get-deps update-deps

GOPATH := $(BUILD_PATH)

SOURCE=$(shell find . -name '*.go')

all: get-deps $(BIN)

clean:
	rm -rf $(BUILD_PATH) $(BIN)

get-deps: .gobuild

.gobuild:
	mkdir -p $(C0_PATH)
	cd "$(C0_PATH)" && ln -s ../../../.. $(PROJECT)

	#
	# Fetch public dependencies via `go get`
	GOPATH=$(GOPATH) go get -d -v github.com/catalyst-zero/$(PROJECT)

$(BIN): $(SOURCE)
	GOPATH=$(GOPATH) go build -o $(BIN)

run-tests:
	GOPATH=$(GOPATH) go test ./...

fmt:
	gofmt -l -w .
