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
	@docker build -t deployments_builder:latest		-f $(GOBASE)/build/package/builder/Dockerfile .
	@docker build -t deployments_cache:latest 		-f $(GOBASE)/build/package/cache/Dockerfile .
	@docker build -t deployments_previewer:latest	-f $(GOBASE)/build/package/previewer/Dockerfile .
	@docker build -t deployments_proxy:latest		-f $(GOBASE)/build/package/proxy/Dockerfile .

.PHONY: run
run:
	@echo "  >  Starting microservices"
	@docker-compose -f deployments/docker-compose.yml up -d

.PHONY: test
test:
	@echo "  >  Making integration tests"
	@set -e ; \
	docker-compose -f deployments/docker-compose.test.yml up --build -d ; \
	sleep 10 ; \
	exitCode=0 ; \
	docker-compose -f deployments/docker-compose.test.yml \
		run -e CGO_ENABLED=0 -e GOOS=linux integration_tests go test || exitCode=$$? ; \
	docker-compose -f deployments/docker-compose.test.yml down ; \
	exit $$exitCode

.PHONY: down
down:
	@echo "  >  Stopping microservices"
	@docker-compose -f deployments/docker-compose.yml down

.PHONY: clean
clean:
	@echo "  >  Cleaning microservice Docker images"
	@IMAGES="$(shell docker images --filter=reference='*deployments_*' -q)"; docker rmi $$IMAGES

.PHONY: ci-build
ci-build:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=mod -o $(GOBIN)/cache		$(GOBASE)/cmd/cache/*.go
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=mod -o $(GOBIN)/previewer	$(GOBASE)/cmd/previewer/*.go
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=mod -o $(GOBIN)/proxy		$(GOBASE)/cmd/proxy/*.go

.PHONY: ci-test
ci-test:
	@go test -mod=mod -race -count 100 $(GOBASE)/internal/domain/usecases/...

.PHONY: ci-lint
ci-lint:
	@golangci-lint run ./...

.PHONY: ci-clean
ci-clean:
	@-rm -fR $(GOBIN)
    @GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
