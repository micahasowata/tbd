.PHONY: run
run:
	go run ./cmd/api

.PHONY: test
test:
	GOFLAGS="-count=1" go test -v -cover -race ./...
