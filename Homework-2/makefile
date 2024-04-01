ifeq ($(POSTGRES_SETUP),)
	POSTGRES_SETUP := user=postgres password=password dbname=postgres host=localhost port=5433 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql 

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down

.PHONY: docker-compose-up
docker-compose-up: docker-compose-build
	docker-compose up

.PHONY: docker-compose-build
docker-compose-build:
	docker-compose build