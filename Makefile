.PHONY: run
run:
	go run ./cmd/api

.PHONY: test
test:
	go test -v -race -buildvcs ./...
	
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./... 
	go tool cover -html=/tmp/coverage.out
