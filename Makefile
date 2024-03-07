# Enhanced Makefile for building Go projects for Linux, macOS (amd64 and arm64)
# with version tagging, and packaging them with a checksum.

.PHONY: all linux mac mac-arm clean

APPNAME := ip-kill
VERSION := $(shell git describe --tags --always)
BUILD_DIR := dist
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"
GOBUILD := go build $(LDFLAGS)

all: linux mac mac-arm

linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/linux-amd64/$(APPNAME)
	cd $(BUILD_DIR)/linux-amd64 && sha256sum $(APPNAME) > $(APPNAME).sha256
	cp LICENSE $(BUILD_DIR)/linux-amd64/
	cd $(BUILD_DIR) && tar czvf $(APPNAME)-linux-amd64-$(VERSION).tar.gz linux-amd64

mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/mac-amd64/$(APPNAME)
	cd $(BUILD_DIR)/mac-amd64 && sha256sum $(APPNAME) > $(APPNAME).sha256
	cp LICENSE $(BUILD_DIR)/mac-amd64/
	cd $(BUILD_DIR) && tar czvf $(APPNAME)-mac-amd64-$(VERSION).tar.gz mac-amd64

mac-arm:
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/mac-arm64/$(APPNAME)
	cd $(BUILD_DIR)/mac-arm64 && sha256sum $(APPNAME) > $(APPNAME).sha256
	cp LICENSE $(BUILD_DIR)/mac-arm64/
	cd $(BUILD_DIR) && tar czvf $(APPNAME)-mac-arm64-$(VERSION).tar.gz mac-arm64

clean:
	rm -rf $(BUILD_DIR)
