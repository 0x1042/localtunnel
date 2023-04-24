export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

.PHONY: fmt build vet clean

clean:
	rm -fr bin

vet:
	go vet ./...

fmt:
	go mod tidy
	gofumpt -l -w .

build:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/lt

build_linux:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/lt

