.PHONY: fmt
fmt: ## format basecode using cargo.
	cargo fmt

.PHONY: lint-fmt
lint-fmt: ## checkc that code format is ok.
	cargo fmt --all --check

.PHONY: lint-clippy
lint-clippy: ## apply linter clippy.
	cargo clippy  --future-incompat-report -- -D warnings

.PHONY: test
test: ## run unit tests using cargo.
	cargo test

.PHONY: run-postgres
run-postgres:
	docker run --name some-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres

.PHONY: start-services
start-services: ## start database
	podman-compose -f docker-compose.yml up

.PHONY: stop-services
stop-services: ## stop database
	podman-compose --file docker-compose.yml down --volumes

.PHONY: connect-db
connect-db: ## connect to postgresql database
	psql -U kbdb -h localhost -p 5432

.PHONY: run-migration
run-migration: ## run pending migrations
	sqlx migrate run --database-url postgres://kbdb:kbpwd@localhost:5432/kbdb

.PHONY: test-get-kb-by-id
test-get-kb-by-id : ## test get kb by id in the persistence layer.
	INTEGRATION_TEST=1 cargo test --package kbservice --bin kbservice \
	-- adapters::pgstorage::pgdb_test::pgstorage_test::test_get_kb_by_id \
	--exact --show-output