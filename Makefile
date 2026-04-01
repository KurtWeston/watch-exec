.PHONY: build install clean test

build:
	go build -o watch-exec .

install:
	go install .

clean:
	rm -f watch-exec

test:
	go test -v ./...

run:
	go run . $(ARGS)
