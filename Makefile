VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BINARY = hanimo

LDFLAGS = -ldflags "-s -w -X main.version=$(VERSION) \
	-X 'github.com/flykimjiwon/hanimo/internal/config.DebugMode=true'"

.PHONY: build run clean build-all install lint test build-index

# Knowledge index
build-index:
	@echo "Building knowledge index..."
	@go run ./cmd/build-index/
	@echo "Index built."

build: build-index
	go build $(LDFLAGS) -o $(BINARY) ./cmd/hanimo

run: build
	./$(BINARY)

install:
	go install $(LDFLAGS) ./cmd/hanimo

clean:
	rm -f $(BINARY)
	rm -rf dist/

build-all: clean
	mkdir -p dist
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64 ./cmd/hanimo
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64 ./cmd/hanimo
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe ./cmd/hanimo
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 ./cmd/hanimo
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 ./cmd/hanimo

lint:
	go vet ./...

test:
	go test ./...
