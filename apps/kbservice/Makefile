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