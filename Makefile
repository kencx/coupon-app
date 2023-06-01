.PHONY: help test dcb dcu dcud seed

default: help

help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## test: run all tests
test:
	go test -v

## dcb: build Dockerfile
dcb:
	docker compose build app

## dcu: bring up prod docker stack
dcu:
	docker compose up -d

## dcud: bring up dev docker stack
dcud:
	docker compose --file docker-compose.dev.yml up -d

## seed: seed database
seed: dev
	curl -s -X POST --header 'application/json' -d @payload.json http://localhost:8080/api/coupons
