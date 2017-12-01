PKGS ?= $(shell glide novendor)
PKG_FILES ?= *.go

# Update path variable to have .bin directory (protoc-gen-go)
export PATH := ./.bin:$(PATH)

default: test

# Installs dependencies
.PHONY: dependencies
dependencies:
	@echo "Installing Glide and locked dependencies..."
	glide --version || go get -u -f github.com/Masterminds/glide
	glide install

# Cleans up produced artifacts
.PHONY: clean
	rm -rf ./gen

# Cleans up dependencies and produced artifacts
.PHONY: purge
purge:
	make clean
	rm -rf ./vendor
	rm -rf ./bin

.PHONY: test
test: .gen/hello.pb.go
	@.build/test.sh

# Protocol buffers generator
.gen/hello.pb.go: .gen proto/hello.proto .bin/protoc-gen-go
	protoc --go_out=.gen proto/*.proto

.gen:
	mkdir .gen

# Builds protobuf generator plugin
.bin/protoc-gen-go:
	./.build/build_vendored.sh .bin github.com/golang/protobuf/protoc-gen-go

