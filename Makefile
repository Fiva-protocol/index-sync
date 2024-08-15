ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: update
update:
	go mod tidy
	go mod verify

.PHONY: build
build:
	sam build

.PHONY: run-local
run-local: build
	sam local start-api

.PHONY: deploy-testnet
deploy-testnet: build
	sam deploy --config-env testnet