binary = coupon-app

.PHONY: help build run clean cover unit-test test dcu

default: help

help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: build binary
build:
	go build -v -o ${binary} .

## run: run binary
run:
	./${binary}

## clean: remove binaries, dist
clean:
	if [ -f ${binary} ]; then rm ${binary}; fi
	go clean

## cover: get code coverage
cover:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out

## test: run all tests
test:
	go test -v

## dcb: build Dockerfile
dcb:
	docker compose build app
