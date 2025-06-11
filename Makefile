include .env

APP_NAME=go_article_api_app

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" up

.PHONY: db/migrations/up/test
db/migrations/up/test: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database "postgres://${TEST_DB_USER}:${TEST_DB_PASSWORD}@localhost:${TEST_DB_PORT}/${TEST_DB_NAME}?sslmode=disable" up

.PHONY: test
test:
	@echo 'Running testings...'
	docker exec $(APP_NAME) go test -v ./internal/tests/...
