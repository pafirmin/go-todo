include .envrc

# ================================================================ #
# HELPERS
# ================================================================ #

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ================================================================ #
# DEVELOPMENT
# ================================================================ #

.PHONY: run/api
run/api:
	go run ./cmd/app -db-address=${DB_ADDR}

.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./db/migrations -database ${DB_ADDR} up

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./db/migrations ${name}

