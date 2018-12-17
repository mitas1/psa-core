GO_SRC = $(shell find . -name '*.go' -type f)

run:
	make -C core run

