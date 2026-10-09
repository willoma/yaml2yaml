# Also update the version in .github/workflows/lint.yml
GOLANGCILINT_VERSION=2.14.0

.PHONY: build check lint lintfix test vulncheck tidy

build: dist/yaml2yaml

check: lint test vulncheck

lint: dist/golangci-lint
	dist/golangci-lint run

lintfix: dist/golangci-lint
	dist/golangci-lint run --fix

test:
	go tool ginkgo

vulncheck:
	go tool govulncheck ./...

tidy:
	go mod tidy

dist/golangci-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b dist v$(GOLANGCILINT_VERSION)

.SECONDEXPANSION:

dist/yaml2yaml:  $$(addsuffix /*.go,$$(filter $(CURDIR)%,$$(shell go list -deps -f '{{.Dir}}' ./cmd/yaml2yaml)))
	CGO_ENABLED=0 go build -o dist/yaml2yaml cmd/yaml2yaml/main.go
