.PHONY: help build build_local up down logs ps test migrate dry-migrate generate
.DEFAULT_GOAL := help

DOCKER_TAG := latest

build: ## Build docker image to deploy
			docker build -t ito0804takuya/gotodo:${DOCKER_TAG} --target deploy ./
build_local: ## Build docker image to local development
			docker compose build --no-cache
up: ## Do docker compose up with hot relead
			docker compose up -d
down: ## Do docker compose down
			docker compose down
logs: ## Tail docker compose logs
			docker compose logs -f
ps: ## Check container status
			docker compose ps
test: ## Execute tests
			go test -race -shuffle=on ./...
migrate: ## Migrate DB
			mysqldef -u todo -p todo -h 127.0.0.1 -P 33306 todo < ./_tools/mysql/schema.sql
dry-migrate: ## Dry Migrate DB
			mysqldef -u todo -p todo -h 127.0.0.1 -P 33306 todo --dry-run < ./_tools/mysql/schema.sql
generate: ## Generate Code
			go generate ./...
help: ## Show options
			@grep -E '^[a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | \
					awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'