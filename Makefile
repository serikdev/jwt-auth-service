
APP_NAME=jwt-service
DOCKER_COMPOSE=docker compose -f deployment/docker-compose.yml



build:
	go build -o $(APP_NAME) main.go


run:
	go run cmd/api/main.go


docker-up:
	$(DOCKER_COMPOSE) up --build


docker-down:
	$(DOCKER_COMPOSE) down


migrate-up:
	bash scripts/migrate-up.sh

migrate-down:
	bash scripts/migrate-down.sh


