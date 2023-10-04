export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

TAGS=urfave_cli_no_docs,netgo
BUILD=go build -tags $(TAGS) -trimpath
LDFLAGS := -s -w -X main.Version="v0.2.4" -X main.Date=$(shell date +"%Y-%m-%d")
TARGET := bin/localtunnel

.PHONY: fmt build vet clean

clean:
	rm -fr bin

vet:
	go vet ./...

fmt:
	go install mvdan.cc/gofumpt@latest
	go mod tidy
	gofumpt -l -w .

build: clean fmt vet
	env CGO_ENABLED=0 $(BUILD) -ldflags "$(LDFLAGS)" -o "$(TARGET)"

build_linux:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD) -ldflags "$(LDFLAGS)" -o "$(TARGET)"

lint:
	golangci-lint run -v
