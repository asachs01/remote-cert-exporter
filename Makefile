.PHONY: build test docker clean

build:
	go build -o cert-exporter

test:
	go test -v ./...

docker:
	docker build -t cert-exporter .

clean:
	rm -f cert-exporter

lint:
	golangci-lint run

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.DEFAULT_GOAL := build 