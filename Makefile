.PHONY: vendors
vendors:
	go mod download

.PHONY: test
test: lint unit

.PHONY: unit
unit:
	go test -race -cover ./...

.PHONY: lint
lint:
	golangci-lint run
	golint -set_exit_status ./...
