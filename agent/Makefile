BINARY=clawwatch-agent
GO=go

.PHONY: build test clean

build:
	$(GO) build -o $(BINARY) .

test:
	$(GO) test ./...

clean:
	rm -f $(BINARY)
