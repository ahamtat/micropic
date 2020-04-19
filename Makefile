PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/bin

.PHONY: gen
gen:
	@echo "  >  Generating code on specifications"
	@protoc $(GOBASE)/api/preview.proto --proto_path=$(GOBASE)/api --go_out=plugins=grpc:$(GOBASE)/internal/adapters/grpc
	@protoc $(GOBASE)/api/cache.proto --proto_path=$(GOBASE)/api --go_out=plugins=grpc:$(GOBASE)/internal/adapters/grpc
	@protoc $(GOBASE)/api/previewer.proto --proto_path=$(GOBASE)/api --go_out=plugins=grpc:$(GOBASE)/internal/adapters/grpc
