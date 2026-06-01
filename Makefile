.PHONY: install build clean release

install:
	go build -o $(shell go env GOPATH)/bin/ddx .

build:
	go build -o ddx .

dev:
	go build -o npm/bin/ddx .

clean:
	rm -f ddx

release:
	./scripts/release.sh
