GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
BINARY_NAME=loglint
BINARY_DIR=bin

all: clean deps test build system_install

.PHONY: deps
deps:
	$(GOCMD) get

.PHONY: test
test:
	$(GOTEST) -v ./test/

.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) ./main.go

.PHONY: run
run:
	$(GORUN) ./main.go

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

.PHONY: system_install
system_install:
	$(GOINSTALL)
