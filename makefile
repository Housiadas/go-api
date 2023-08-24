DOCKER_COMPOSE_LOCAL := docker-compose -f ./docker/local/docker-compose.yml
MIGRATE := $(DOCKER_COMPOSE_LOCAL) run --rm app migrate
DB_DSN_LOCAL := postgres://housi:secret123@db/housi_db?sslmode=disable

docker-build:
	$(DOCKER_COMPOSE_LOCAL) build --no-cache --pull

docker-up:
	$(DOCKER_COMPOSE_LOCAL) up -d app

docker-down:
	$(DOCKER_COMPOSE_LOCAL) down --remove-orphans

docker-clean:
	docker system prune && \
    docker image prune && \
    docker volume prune

migrate-create:
	read -p "Enter filename:" cmd; \
	$(MIGRATE) create -seq -ext=.sql -dir=./migrations $$cmd

migrate-up:
	$(MIGRATE) -path=./migrations -database=$(DB_DSN_LOCAL) up

migrate-down:
	$(MIGRATE) -path=./migrations -database=$(DB_DSN_LOCAL) down
