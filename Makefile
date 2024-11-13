# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=smol-docker

DIST_FOLDER=dist

build:
		$(GOBUILD) -o $(BINARY_NAME) -v
clean:
		$(GOCLEAN)
		rm -rf $(DIST_FOLDER)
build-all:
		mkdir -p $(DIST_FOLDER)
		# [darwin/amd64] - Intel Mac
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_darwin_amd64 -v
		# [darwin/arm64] - Apple Silicon Mac
		CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_darwin_arm64 -v
		# [linux/amd64] - 64-bit Linux
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_linux_amd64 -v
		# [linux/386] - 32-bit Linux
		CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GOBUILD) -o $(DIST_FOLDER)/$(BINARY_NAME)_linux_386 -v