##@ Golang

.PHONY: fmt
fmt: ## Run go fmt against code.
	@$(LOG_TARGET)
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	@$(LOG_TARGET)
	go vet ./...

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter & yamllint
	@$(LOG_TARGET)
	$(GOLANGCI_LINT) run

.PHONY: build
build: manifests generate fmt vet ## Build manager binary.
	@$(LOG_TARGET)
	go build -o operator cmd/main.go


.PHONY: compile
compile: fmt vet ## Compile binary
	@$(LOG_TARGET)
	go build -o operator cmd/main.go

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	@$(LOG_TARGET)
	go run ./cmd/main.go