export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

TAGS=urfave_cli_no_docs,netgo
GOEXE=/opt/hostedtoolcache/go/1.22.0/x64/bin/go
BUILD=$(GOEXE) build -tags $(TAGS) -trimpath
VER=$(shell git rev-parse --short HEAD)
DATE=$(shell date +"%Y-%m-%d")
LDFLAGS := -s -w -X main.Version=$(VER) -X main.Date=$(DATE)
TARGET := bin/localtunnel

.PHONY: fmt build vet clean build_local

clean:
	rm -fr bin

vet:
	$(GOEXE) vet ./...

fmt:
	export GOROOT=/opt/hostedtoolcache/go/1.22.0/x64
	$(GOEXE) install mvdan.cc/gofumpt@latest
	$(GOEXE) version
	$(GOEXE) env
	gofumpt -l -w .

build: clean fmt vet
	env CGO_ENABLED=0 $(BUILD) -ldflags "$(LDFLAGS)" -o "$(TARGET)"

build_linux:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD) -ldflags "$(LDFLAGS)" -o "$(TARGET)"

lint:
	golangci-lint run -v

build_local:clean
	gofumpt -l -w .
	go vet ./...
	env CGO_ENABLED=0 go build -tags $(TAGS) -trimpath -ldflags "$(LDFLAGS)" -o "$(TARGET)"