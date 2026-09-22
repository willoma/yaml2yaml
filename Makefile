build: dist/yaml2yaml

lint:
	@mkdir -p $(shell go env GOCACHE) $(shell go env GOMODCACHE) $(wildcard ~/.cache)/golangci-lint
	docker run --rm -t -v $(CURDIR):/app -w /app \
	--user $(shell id -u):$(shell id -g) \
	-v $(shell go env GOCACHE):/.cache/go-build -e GOCACHE=/.cache/go-build \
	-v $(shell go env GOMODCACHE):/.cache/mod -e GOMODCACHE=/.cache/mod \
	-v $(wildcard ~/.cache)/golangci-lint:/.cache/golangci-lint -e GOLANGCI_LINT_CACHE=/.cache/golangci-lint \
	golangci/golangci-lint:v2.13.2 golangci-lint run

lintfix:
	@mkdir -p $(shell go env GOCACHE) $(shell go env GOMODCACHE) $(wildcard ~/.cache)/golangci-lint
	docker run --rm -t -v $(CURDIR):/app -w /app \
	--user $(shell id -u):$(shell id -g) \
	-v $(shell go env GOCACHE):/.cache/go-build -e GOCACHE=/.cache/go-build \
	-v $(shell go env GOMODCACHE):/.cache/mod -e GOMODCACHE=/.cache/mod \
	-v $(wildcard ~/.cache)/golangci-lint:/.cache/golangci-lint -e GOLANGCI_LINT_CACHE=/.cache/golangci-lint \
	golangci/golangci-lint:v2.13.2 golangci-lint run --fix

dist/yaml2yaml:
	CGO_ENABLED=0 go build -o dist/yaml2yaml cmd/yaml2yaml/main.go
