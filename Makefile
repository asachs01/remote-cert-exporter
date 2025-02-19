.PHONY: build test docker clean

build:
	go build -o remote-cert-exporter

test:
	go test -v ./...

docker:
	docker build -t remote-cert-exporter .

clean:
	rm -f remote-cert-exporter

lint:
	golangci-lint run

format:
	go fmt ./...

check-format:
	@if [ -n "$$(go fmt ./...)" ]; then \
		echo "Files not properly formatted. Run 'make format'"; \
		exit 1; \
	fi

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.DEFAULT_GOAL := build 