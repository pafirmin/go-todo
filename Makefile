# ================================================================ #
# HELPERS
# ================================================================ #

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh ${REMOTE_USER}@${PROD_HOST_IP}

## production/deploy: deploy the api to production
.PHONY: production/deploy
production/deploy:
	rsync -P ./bin/linux_amd64/app ${REMOTE_USER}@${PROD_HOST_IP}:${REMOTE_PATH}
	rsync -rP --delete ./db/migrations ${REMOTE_USER}@${PROD_HOST_IP}:${REMOTE_PATH}
	rsync -P ./remote/production/go-todo.service ${REMOTE_USER}@${PROD_HOST_IP}:${REMOTE_PATH}
	ssh -t ${REMOTE_USER}@${PROD_HOST_IP} '\
		sudo mv ${REMOTE_PATH}/go-todo.service /etc/systemd/system/ \
		&& sudo systemctl enable go-todo \
		&& sudo systemctl restart go-todo \
	'

.PHONY: production/migrate
production/migrate:
	ssh -t ${REMOTE_USER}@${PROD_HOST_IP} 'migrate -path ${REMOTE_PATH}/migrations -database $$GO_TODO_DB_ADDR up'


# ================================================================ #
# DEVELOPMENT
# ================================================================ #

## run/app: run the application
.PHONY: run/app
run/app:
	go run ./cmd/app -db-address=${DEV_DB_ADDR} -jwt-secret=${JWT_SECRET}

## run/bin: execute application binary
.PHONY: run/bin
run/bin:
	./bin/app -db-address=${DEV_DB_ADDR} -jwt-secret=${JWT_SECRET}

## db/migrations/up: run all database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./db/migrations -database ${DEV_DB_ADDR} up

## db/migrations/new: create a new pair of blank migration files
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./db/migrations ${name}

# ================================================================ #
# BUILD
# ================================================================ #

## build/api: build the application
.PHONY: build/api
build/api:
	@echo 'Building cmd/app...'
	go build -ldflags='-s' -o=./bin/app ./cmd/app
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/app ./cmd/app

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and vet, format and test code
.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Tidying and verifying dependencies...'
	go mod tidy
	go mod verify
	@echo 'Running tests...'
	go test -race -vet=off ./...
