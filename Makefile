export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w -X main.Version="v0.2.2" -X main.Date=$(shell date +"%Y-%m-%d")
TARGET := bin/lt

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
	env CGO_ENABLED=0 go build -tags urfave_cli_no_docs,netgo -trimpath -ldflags "$(LDFLAGS)" -o "$(TARGET)"

build_linux:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o "$(TARGET)"

