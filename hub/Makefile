BINARY=clawwatch-hub
AGENT_BINARY=clawwatch-agent
GO=go

.PHONY: build build-agent build-linux run clean

build:
	$(GO) build -o $(BINARY) .

build-agent:
	cd ../agent && $(GO) build -o ../hub/$(AGENT_BINARY) .

build-linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GO) build -o $(BINARY)-linux-arm64 .

run: build
	./$(BINARY) --port 4848

run-secure: build
	./$(BINARY) --port 4848 \
		--agent-token $${HUB_AGENT_TOKEN:-change-me} \
		--console-token $${HUB_CONSOLE_TOKEN:-change-me}

clean:
	rm -f $(BINARY) $(BINARY)-linux-* $(AGENT_BINARY)
