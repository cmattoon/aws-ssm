FILES=$(shell go list ./... | grep -v "/vendor/")
build:
	@go build $(FILES)
