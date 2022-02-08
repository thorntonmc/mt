.DEFAULT_GOAL := build

fmt:
	go fmt ./...
PHONY:fmt

vet: fmt
	go vet ./...
.PHONEY:vet

test: vet
	go test ./...

build: test
	go build main.go

clean:
	rm -rf build

