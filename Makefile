SHELL := /bin/bash

TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)

VERSION := 0.1
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean proto

all: proto build

# Using packr (https://github.com/gobuffalo/packr) to build
$(TARGET): $(SRC)
	@packr build $(LDFLAGS) -o $(TARGET)

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET) history/*.pb.go backend/*_pb2*

proto:
	protoc -I=history/ --go_out=plugins=grpc:history/ history/history.proto
	python3 -m grpc_tools.protoc -Ihistory/ --python_out=backend/ --grpc_python_out=backend/ history/history.proto 