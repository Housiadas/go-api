DOCKER_COMPOSE_LOCAL := docker-compose -f ./docker/local/docker-compose.yml

docker-build:
	$(DOCKER_COMPOSE_LOCAL) build --no-cache --pull

docker-up:
	$(DOCKER_COMPOSE_LOCAL) up -d

docker-down:
	$(DOCKER_COMPOSE_LOCAL) down

docker-clean:
	docker system prune && \
    docker image prune && \
    docker volume prune

local-up:
	go run ./cmd/api -port=3030 -env=production
