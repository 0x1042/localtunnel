export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

TAGS=urfave_cli_no_docs,netgo
BUILD=/opt/hostedtoolcache/go/1.21.2/x64/bin/go build -tags $(TAGS) -trimpath
LDFLAGS := -s -w -X main.Version="v0.2.4" -X main.Date=$(shell date +"%Y-%m-%d")
TARGET := bin/localtunnel

.PHONY: fmt build vet clean

clean:
	rm -fr bin

vet:
	/opt/hostedtoolcache/go/1.21.2/x64/bin/go vet ./...

fmt:
	export GOROOT=/opt/hostedtoolcache/go/1.21.2/x64
	/opt/hostedtoolcache/go/1.21.2/x64/bin/go install mvdan.cc/gofumpt@latest
	/opt/hostedtoolcache/go/1.21.2/x64/bin/go version
	/opt/hostedtoolcache/go/1.21.2/x64/bin/go env
	/opt/hostedtoolcache/go/1.21.2/x64/bin/go mod tidy -v -go=1.21
	gofumpt -l -w .

build: clean fmt vet
	env CGO_ENABLED=0 $(BUILD) -ldflags "$(LDFLAGS)" -o "$(TARGET)"

build_linux:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD) -ldflags "$(LDFLAGS)" -o "$(TARGET)"

lint:
	golangci-lint run -v
