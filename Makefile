.PHONY: all build run test clean

BINARY_NAME=hyena-ufa-plugin

all: build

build:
	go build -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME) $(ARGS)

test:
	go test -v ./...

clean:
	go clean; rm -f $(BINARY_NAME)
