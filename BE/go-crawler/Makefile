SHELL := /bin/bash
GOLANG_VERSION := 1.9
APP_NAME := 'go-crawler'
BRANCH_NAME ?= $(shell git name-rev --name-only HEAD | cut -d '/' -f3)
USERNAME := $(shell whoami)
USERID := $(shell id -u)
GROUPNAME := $(shell id -gn)
GROUPID := $(shell id -g)
VERSION ?= 0.0.1
BUILD_NUM := $(if $(BUILD_NUMBER),$(BUILD_NUMBER),0)
APP_HASH := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags="-s -X main.ApplicationVersion=${VERSION}"

default: build
.PHONY: default

help:
	@echo "  build              build app sources"
	@echo "  install            installs the packages named by the import paths"
	@echo "  lint               run lint tests"
	@echo "  clean              clean-up"
.PHONY: help

build: clean
	@mkdir -p ./apps/$(APP_NAME)
	@export GOPATH="$(CURDIR)" && pushd ./src/$(APP_NAME)/ && go build $(LDFLAGS) -a -o ../../apps/$(APP_NAME)/$(APP_NAME) && popd
.PHONY: build

install:
	@export GOPATH="$(CURDIR)" && pushd ./src/$(APP_NAME)/ && go install -v && popd
.PHONY: install

lint: install
	@export GOPATH="$(CURDIR)" && pushd ./src/$(APP_NAME)/ && gometalinter --deadline=180s --vendor ./... && popd
.PHONY: lint

vendor:
	export GOPATH="$(CURDIR)" && pushd ./src/ && govendor sync && govendor fetch +outside && popd
.PHONY: vendor

clean:
	@rm -Rf ./apps ./pkg ./src/$(APP_NAME)/$(APP_NAME)
.PHONY: clean

