
PROG =	bodyworkplans

all: run
	./${PROG}

run: build

build:
	go build -o ${PROG}

gen: generate

generate:
	go generate ./...

fmt:
	gofmt -w $$(find . -name "*.go" | grep -v -e data/exercises.go -e data/plans.go)

bundle:
	fyne bundle favicon.svg > favicon.go
