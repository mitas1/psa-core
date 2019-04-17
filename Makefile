GO_SRC = $(shell find .. -name '*.go' -type f)
TARGET := psa-core

.DEFAULT_GOAL := all
.PHONY: build run

$(TARGET): $(GO_SRC)
	go build -o $(TARGET)

all: build

build: $(TARGET)

run: build
	./$(TARGET)
