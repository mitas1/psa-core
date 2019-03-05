GO_SRC = $(shell find .. -name '*.go' -type f)
TARGET := psa-core

# If the first argument is "run"...
ifeq (run,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

.DEFAULT_GOAL := all
.PHONY: build run

$(TARGET): $(GO_SRC)
	go build -o $(TARGET)

all: build

build: $(TARGET)

run: build
	@echo ./$(TARGET) $(RUN_ARGS)
	./$(TARGET) $(RUN_ARGS)
