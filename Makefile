include .env

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

.PHONY: db
db:
	echo ${DB_DSN}
	echo ${TEST_DB_DSN}

.PHONY: db/migrations/new
db/migrations/new:
	migrate create -seq -ext .sql -dir migrations ${name}

.PHONY: db/migrations/up
db/migrations/up:
	migrate -path migrations -database ${DB_DSN} up
	migrate -path migrations -database ${TEST_DB_DSN} up

.PHONY: db/migrations/down 
db/migrations/down:
	migrate -path migrations -database ${DB_DSN} down
	migrate -path migrations -database ${TEST_DB_DSN} down

.PHONY: db/migrations/force 
db/migrations/force:
	migrate -path migrations -database ${DB_DSN} force ${v}
	migrate -path migrations -database ${TEST_DB_DSN} force ${v}

.PHONY: db/migrations/drop 
db/migrations/drop:
	migrate -path migrations -database ${DB_DSN} drop
	migrate -path migrations -database ${TEST_DB_DSN} drop
