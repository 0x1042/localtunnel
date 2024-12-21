export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

TAGS=urfave_cli_no_docs,netgo
GOEXE=go
BUILD=$(GOEXE) build -tags $(TAGS) -trimpath
VER=$(shell git rev-parse --short HEAD)
DATE=$(shell date +"%Y-%m-%d")
LDFLAGS := -s -w -X main.Version=$(VER) -X main.Date=$(DATE)
TARGET := bin/localtunnel

.PHONY: fmt build vet clean build_local info

info:
	@echo "go env"
	go env

clean:
	rm -fr bin

vet:
	$(GOEXE) vet ./...

fmt:
	$(GOEXE) install mvdan.cc/gofumpt@latest
	$(GOEXE) version
	$(GOEXE) env
	gofumpt -l -w .

build: info clean fmt vet
	env CGO_ENABLED=0 $(BUILD) -ldflags "$(LDFLAGS)" -o "$(TARGET)"

build_linux:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BUILD) -ldflags "$(LDFLAGS)" -o "$(TARGET)"

lint:
	golangci-lint run -v

build_local:clean
	gofumpt -l -w .
	go vet ./...
	env CGO_ENABLED=0 go build -tags $(TAGS) -trimpath -ldflags "$(LDFLAGS)" -o "$(TARGET)"

update:
	go get -u ./...
	go mod tidy
