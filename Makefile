GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

help:
	@echo "This is a helper makefile for oapi-codegen"
	@echo "Targets:"
	@echo "    generate:    regenerate all generated files"
	@echo "    test:        run all tests"
	@echo "    tidy         tidy go mod"
	@echo "    lint         lint the project"

$(GOBIN)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.64.5

.PHONY: tools
tools: $(GOBIN)/golangci-lint

lint: tools
	$(GOBIN)/golangci-lint run ./...

lint-ci: tools
	$(GOBIN)/golangci-lint run ./... --out-format=colored-line-number --timeout=5m

generate:
	go generate ./...
test:
	go test -cover ./...

tidy:
	go mod tidy

tidy-ci:
	tidied -verbose
