.PHONY: all build run test clean deploy-dev-binaries

BINARY_NAME=hyena-ufa-plugin

# local-representative's plug-in dev repo watcher passes DEV_PLUGIN_BIN (where
# to build the binary), DEV_PLUGIN_DIR and PLUGIN_APP_ID on the command line.
DEV_PLUGIN_BIN ?= $(CURDIR)/$(BINARY_NAME)

all: build

build:
	go build -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME) $(ARGS)

test:
	go test -v ./...

# The hook local-representative runs when a tagged commit lands (or a
# .ufabuild request is present): test, then build the plug-in to
# $(DEV_PLUGIN_BIN) and nothing more. LR deploys the binary, launches it with
# --lr-host/--lr-port/--port/--name (and --dev-mode) and restarts it after
# each successful build.
deploy-dev-binaries:
	go test ./...
	mkdir -p $(dir $(DEV_PLUGIN_BIN))
	go build -o $(DEV_PLUGIN_BIN) .

clean:
	go clean; rm -f $(BINARY_NAME)
