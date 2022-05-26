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

.PHONY: run/app
run/api:
	go run ./cmd/app -db-address=${DB_ADDR} -jwt-secret=${JWT_SECRET}

.PHONY: run/bin
run/bin:
	./bin/app -db-address=${DB_ADDR} -jwt-secret=${JWT_SECRET}

.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./db/migrations -database ${DB_ADDR} up

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./db/migrations ${name}

# ================================================================ #
# BUILD
# ================================================================ #

.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/app ./cmd/app
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/app ./cmd/app

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Tidying and verifying dependencies...'
	go mod tidy
	go mod verify
	@echo 'Running tests...'
	go test -race -vet=off ./...
