include .env
export

export POSTGRES_URL ?= postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSLMODE}


run:
	docker compose up -d
	go run cli/httpserver/main.go


populate:
	docker compose up -d
	go run cli/populator/main.go