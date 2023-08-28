# Include variables from the local .env file
include ./docker/local/.env

DOCKER_COMPOSE_LOCAL := docker-compose -f ./docker/local/docker-compose.yml
MIGRATE := $(DOCKER_COMPOSE_LOCAL) run --rm app migrate

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## docker/build: Build all the containers
.PHONY: docker/build
docker/build:
	$(DOCKER_COMPOSE_LOCAL) build --no-cache --pull

## docker/up: Start all the containers for the app
.PHONY: docker/up
docker/up:
	$(DOCKER_COMPOSE_LOCAL) up -d app

## docker/stop: stop all containers
.PHONY: docker/stop
docker/stop:
	$(DOCKER_COMPOSE_LOCAL) stop

## docker/down: stop and remove all containers
.PHONY: docker/down
docker/down:
	$(DOCKER_COMPOSE_LOCAL) down --remove-orphans

## docker/clean: docker clean all
.PHONY: docker/clean
docker/clean:
	docker system prune && \
    docker image prune && \
    docker volume prune

## db/migrations/create name=$1: create new migration files
.PHONY: db/migrations/create
db/migrate/create:
	@echo 'Creating migration files for ${name}...'
	$(MIGRATE) create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrate/up
db/migrate/up: confirm
	$(MIGRATE) -path=./migrations -database=${DB_DSN} up

## db/migrations/down: apply all down database migrations (DROP Database)
.PHONY: db/migrate/down
db/migrate/down: confirm
	$(MIGRATE) -path=./migrations -database=${DB_DSN} down

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api
