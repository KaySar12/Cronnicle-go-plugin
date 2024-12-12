GOCMD=go
GOCLEAN=$(GOCMD) clean
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
GOBUILD=$(GOCMD) build
BUILD_PATH = $(BASE_PATH)/build
BASE_PATH := $(shell pwd)
APP_NAME=ndutils
MAIN= $(BASE_PATH)/main.go
MASSDNS_PATH= $(BASE_PATH)/massdns
build_bin:
	mkdir -p build 
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 $(GOBUILD) -o $(BUILD_PATH)/$(APP_NAME) $(MAIN)

clean_assets:
	rm $(BUILD_PATH)
