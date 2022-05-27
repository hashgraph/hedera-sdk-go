build: format lint
  go build

test-unit: build
  go test -tags="unit" -v

test-e2e: build
  go test -tags="e2e" -v -timeout 9999s

format:
  go fmt || goimports -w *.go

lint:
  golangci-lint run

set-previewnet:
  echo export CONFIG_FILE=""

set-testnet:
  echo export CONFIG_FILE=""
