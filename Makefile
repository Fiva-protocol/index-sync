ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: build
build:
	sam build

.PHONY: run-local
run-local: build
	sam local start-api

.PHONY: deploy
deploy: build
	sam deploy --config-env ${ENV}