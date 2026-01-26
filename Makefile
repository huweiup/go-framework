.PHONY: build test lint clean

build:
	go build -o bin/goframework cmd/goframework/main.go

test:
	go test -v ./pkg/... ./cmd/...

lint:
	go vet ./...

clean:
	rm -rf bin/
