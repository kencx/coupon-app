binary = coupon-app

.PHONY: help build run clean cover test dcb dev seed

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

dev:
	@docker compose up -d db app
	cd ui && docker compose up -d

## seed: seed database
seed: dev
	curl -s -X POST --header 'application/json' -d @payload.json http://localhost:8080/api/coupons
