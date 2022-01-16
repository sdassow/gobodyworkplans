
PROG =	bodyworkplans

all: run
	./${PROG}

run: build

build:
	go build -o ${PROG}

gen: generate

generate:
	go generate ./...
