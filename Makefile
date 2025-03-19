.PHONY: setup test build clean tidy lint fmt

VERSION := $(shell cat VERSION)
LDFLAGS := -ldflags "-X github.com/kreasimaju/auth.Version=$(VERSION)"

setup:
	go mod download

test:
	go test -v ./...

test-cover:
	go test -cover ./...

build:
	go build $(LDFLAGS) -o bin/auth ./...

clean:
	rm -rf bin

tidy:
	go mod tidy

lint:
	golangci-lint run

fmt:
	go fmt ./...

tag:
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)

release: test build tag
	@echo "Released v$(VERSION)" 