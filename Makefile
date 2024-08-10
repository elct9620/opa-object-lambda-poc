.PHONY: build

deploy: build
	sam deploy

build:
	sam build
