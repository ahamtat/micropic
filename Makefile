PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/build/bin

.PHONY: gen
gen:
	@echo "  >  Generating code on specifications"
	@protoc $(GOBASE)/api/preview.proto --proto_path=$(GOBASE)/api --go_out=plugins=grpc:$(GOBASE)/internal/adapters/grpcapi
	@protoc $(GOBASE)/api/cache.proto --proto_path=$(GOBASE)/api --go_out=plugins=grpc:$(GOBASE)/internal/adapters/grpcapi
	@protoc $(GOBASE)/api/previewer.proto --proto_path=$(GOBASE)/api --go_out=plugins=grpc:$(GOBASE)/internal/adapters/grpcapi

.PHONY: build
build:
	@echo "  >  Building microservices binaries & Docker images"
	@env GOOS=linux GOARCH=amd64 go build -mod=mod -o $(GOBIN)/cache        $(GOBASE)/cmd/cache/*.go
	@env GOOS=linux GOARCH=amd64 go build -mod=mod -o $(GOBIN)/previewer	$(GOBASE)/cmd/previewer/*.go
	@env GOOS=linux GOARCH=amd64 go build -mod=mod -o $(GOBIN)/proxy		$(GOBASE)/cmd/proxy/*.go
	@docker build -t deployments_cache:latest 		-f $(GOBASE)/build/package/cache/Dockerfile .
	@docker build -t deployments_previewer:latest	-f $(GOBASE)/build/package/previewer/Dockerfile .
	@docker build -t deployments_proxy:latest		-f $(GOBASE)/build/package/proxy/Dockerfile .

.PHONY: clean
clean:
	@echo "  >  Cleaning build cache"
	@rm -fR $(GOBIN)
	@rm -fR $(GOBASE)/logs
    @env GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean