BIN_OUT=bin
BIN_NAME=parallel
PLATFORMS=windows/386 \
	windows/amd64 \
	darwin/amd64 \
	linux/386 \
	linux/amd64 \
	linux/arm

.PHONY: default
default: host

.PHONY: all
all: host $(PLATFORMS)

.PHONY: host
host:
	go build -o $(BIN_OUT)/$(BIN_NAME) cmd/parallel.go

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	GOOS=$(@D) GOARCH=$(@F) go build -o $(BIN_OUT)/$(BIN_NAME)-$(@D)-$(@F) cmd/parallel.go

.PHONY: clean
clean:
	rm -rf bin/
