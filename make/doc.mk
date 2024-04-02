##@ Doc

DOCS_DIR ?= $(shell pwd)/docs
$(DOCS_DIR):
	mkdir -p $(DOCS_DIR)

api-docs: crdoc kustomize $(DOCS_DIR)  ## Generate Jaeger Operator API Docs
	@$(LOG_TARGET)
	$(KUSTOMIZE) build config/crd -o $(DOCS_DIR)/crd-output.yaml ;\
	$(CRDOC) -r $(DOCS_DIR)/crd-output.yaml -o $(DOCS_DIR)/api.md ;\
