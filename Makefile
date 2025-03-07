# Image URL to use all building/pushing image targets
IMG ?= configsync/operator:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: build

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests
	go test ./... -v

##@ Build

.PHONY: build
build: fmt vet ## Build the operator binary
	go build -o bin/configsync-operator cmd/operator/main.go

.PHONY: run
run: fmt vet ## Run the operator locally
	go run ./cmd/operator/main.go

.PHONY: docker-build
docker-build: ## Build docker image
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image
	docker push ${IMG}

##@ Deployment

.PHONY: install
install: ## Install CRDs into the K8s cluster
	kubectl apply -f deploy/operator.yaml

.PHONY: uninstall
uninstall: ## Uninstall CRDs from the K8s cluster
	kubectl delete -f deploy/operator.yaml

.PHONY: deploy
deploy: ## Deploy controller in the configured Kubernetes cluster
	kubectl apply -f deploy/operator.yaml

.PHONY: undeploy
undeploy: ## Undeploy controller from the configured Kubernetes cluster
	kubectl delete -f deploy/operator.yaml

##@ Helpers

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help 