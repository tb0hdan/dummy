SHELL=/bin/bash
IMAGE_TAG := $(shell git rev-parse HEAD)
IMAGE_NAME := "quay.io/companyname/dummy_project"

.PHONY: ci
ci: deps deps_check lint build test

.PHONY: deps
deps:
	dep ensure -v

.PHONY: deps_check
deps_check:
	@test -z "$(shell git status -s ./vendor ./Gopkg.*)"

.PHONY: grpcgen
grpcgen:
	go build -o protoc-gen-go ./vendor/github.com/golang/protobuf/protoc-gen-go
	protoc -I ./vendor -I api api/service.proto --plugin=./protoc-gen-go --go_out=plugins=grpc:api
	go build ./api

.PHONY: build
build:
	go build -o artifacts/svc

.PHONY: run
run:
	go run ./main.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -cover -v `go list ./...`

.PHONY: dockerise
dockerise:
	docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" -f Dockerfile .

.PHONY: mock_service_db
mock_service_db:
	mockgen -source=service/storage.go -destination=service/mock/storage.go

