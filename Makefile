GO_SRC = $(shell find . -name '*.go' -type f)

build:
	g++ -o core/check_solution _check/check_solution.cpp

run: build
	make -C core run

