ifeq ($(POSTGRES_SETUP),)
	POSTGRES_SETUP := user=postgres password=password dbname=postgres host=localhost port=5433 sslmode=disable
endif

ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=postgres password=password dbname=postgres host=localhost port=5434 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
MOCKGEN_TAG=v1.6.0
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql 

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: docker-compose-up
docker-compose-up: docker-compose-build
	docker-compose up

.PHONY: docker-compose-build
docker-compose-build:
	docker-compose build

.PHONY: .generate-mockgen-deps
.generate-mockgen-deps:
ifeq ($(wildcard $(MOCKGEN_BIN)),)
	@GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@$(MOCKGEN_TAG)
endif

.PHONY: .generate-mockgen
.generate-mockgen:
	PATH="$(LOCAL_BIN):$$PATH" go generate -x -run=mockgen ./...

.test:
	$(info Running tests...)
	go test ./internal/...

test: .test

test-gen: .generate-mockgen .test

integration:
	go test ./tests --tags=integration
