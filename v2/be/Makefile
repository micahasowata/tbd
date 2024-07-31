include .env
# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## confirm: ensure critical commands are confirmed before running
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## db: display the data source names of all databases connected the application
.PHONY: db
db:
	@echo ${DSN}

## new: create new up and down migration scripts
.PHONY: new
new:
	@migrate create -seq -ext .sql -dir migrations ${name}

## up: run all up migration scripts against all connected databases
.PHONY: up
up: confirm
	@migrate -path migrations -database ${DSN} up

## down: run all down migration scripts against all connected databases
.PHONY: down
down:
	@migrate -path migrations -database ${DSN} down

## force: revert schema structure to the specified version across all connected databases
.PHONY: force
force: confirm
	@migrate -path migrations -database ${DSN} force ${v}

## drop: drop all connected databases
.PHONY: drop
drop: confirm
	@migrate -path migrations -database ${DSN} drop
