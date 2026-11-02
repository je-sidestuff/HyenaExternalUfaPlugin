.PHONY: all build run test clean deploy-dev-binaries

BINARY_NAME=hyena-ufa-plugin

# Where deploy-dev-binaries runs the plug-in from, and its pid/log files.
# LR_PORT and ARGS (e.g. ARGS="--dev-mode") change how it is launched.
DEPLOY_DIR ?= $(CURDIR)/.deploy
LR_PORT ?= 8082

all: build

build:
	go build -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME) $(ARGS)

test:
	go test -v ./...

# The hook local-representative's plug-in dev repo watcher runs when a tagged
# commit lands (or a .ufabuild request is present): test, build, then stop the
# previous instance and relaunch, so the plug-in says hello again.
deploy-dev-binaries:
	go test ./...
	mkdir -p $(DEPLOY_DIR)
	go build -o $(DEPLOY_DIR)/$(BINARY_NAME).new .
	-if [ -f $(DEPLOY_DIR)/pid ]; then kill `cat $(DEPLOY_DIR)/pid` 2>/dev/null; sleep 1; fi
	mv $(DEPLOY_DIR)/$(BINARY_NAME).new $(DEPLOY_DIR)/$(BINARY_NAME)
	(cd $(DEPLOY_DIR) && exec nohup ./$(BINARY_NAME) --lr-port $(LR_PORT) $(ARGS) > plugin.log 2>&1 < /dev/null) & echo $$! > $(DEPLOY_DIR)/pid

clean:
	go clean; rm -f $(BINARY_NAME)
