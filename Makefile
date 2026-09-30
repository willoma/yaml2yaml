# Also update the version in .github/workflows/lint.yml
GOLANGCILINT_VERSION=2.14.0

build: dist/yaml2yaml

lint: dist/golangci-lint
	dist/golangci-lint run

lintfix: dist/golangci-lint
	dist/golangci-lint run --fix

test:
	go tool ginkgo

dist/golangci-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b dist v$(GOLANGCILINT_VERSION)

dist/yaml2yaml: *.go
	CGO_ENABLED=0 go build -o dist/yaml2yaml cmd/yaml2yaml/main.go
