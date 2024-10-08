CONTAINERCMD?=docker
COMMIT_HASH?=$(shell git describe --dirty --tags --always)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION?=0.1.0

.PHONY: run
run: ## run the application using cargo. remember to add KB_DB_PWD
	KB_WEB_SERVER_PORT=3030 \
	KB_DB_NAME=kbdb \
	KB_DB_USER=kbdb \
	KB_DB_HOST=localhost \
	KB_DB_PORT=5432 \
	RUST_LOG=debug cargo run

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

.PHONY: start-db-services
start-db-services: ## start database
	podman-compose -f docker-compose.yml up database

.PHONY: stop-services
stop-services: ## stop database
	podman-compose --file docker-compose.yml down --volumes

.PHONY: connect-db
connect-db: ## connect to postgresql database
	psql -U kbdb -h localhost -p 5432

.PHONY: run-migration
run-migration: ## run pending migrations
	sqlx migrate run --database-url postgres://kbdb:kbpwd@localhost:5432/kbdb

.PHONY: add-kb
add-kb: ## add a new kb
	curl -H "Content-Type: application/json" \
	--data '{"key":"esme", "value":"emerald203", "notes":"her name", "kind":"names", "tags":["name", "emerald"]}' \
	-X POST http://localhost:3030/kbs

.PHONY: search-kbs
get-kbs: ## get the existing kbs in the service
	curl -X GET http://localhost:3030/kbs

.PHONY: search-kbs-with-key
search-kbs-with-key: ## get the existing kbs with a given key pattern in the service
	curl -X GET http://localhost:3030/kbs?key=rick

.PHONY: search-kbs-with-tags
search-kbs-with-tags: ## get the existing kbs with a given tag pattern in the service
	curl -X GET http://localhost:3030/kbs?keyword=names

.PHONY: get-kb-by-id
get-kb-by-id: ## get the existing kb by id in the service
	curl -X GET http://localhost:3030/kbs/6411a28b-640a-43d9-b901-1c4b15d91568

        
# INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS)
# VALUES ('681cca89-890b-4667-8ca0-e328546e268c', 'red', 'remember this color', 'one color', 'concepts', 'color concepts')
# RETURNING KB_ID
.PHONY: test-get-kb-by-id
test-get-kb-by-id : ## test get kb by id in the persistence layer.
	INTEGRATION_TEST=1 cargo test --package kbservice --bin kbservice \
	-- adapters::pgstorage::pgdb_test::pgstorage_test::test_get_kb_by_id \
	--exact --show-output

# INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS)
# VALUES ('681cca89-890b-4667-8ca0-e328546e268c', 'red', 'remember this color', 'one color', 'concepts', 'color concepts')
# RETURNING KB_ID
.PHONY: test-get-kb-by-key
test-get-kb-by-key : ## test get kb by key in the persistence layer.
	INTEGRATION_TEST=1 cargo test --package kbservice --bin kbservice \
	-- adapters::pgstorage::pgdb_test::pgstorage_test::test_get_kb_by_key \
	--exact --show-output

# INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ('6411a28b-640a-43d9-b901-1c4b15d91568', 'frederick', 'long name', 'multiple names', 'names', 'name names');
# INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ('5a2579f7-83b9-4891-8dbc-e0024b5f3505', 'rick', 'short name', 'just one name', 'names', 'name names');
# INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ('22cfc4fb-f9b6-4f6e-9158-9982347ad2a7', 'patrick', 'a saint', 'names', 'words', 'over words');
# SELECT KB_ID, KB_KEY, KIND, TAGS::TEXT AS TAGS FROM kbs WHERE KB_KEY LIKE '%rick%' LIMIT 10 OFFSET 0;
.PHONY: test-search-by-key
test-search-by-key : ## test search by key in the persistence layer.
	INTEGRATION_TEST=1 cargo test --package kbservice --bin kbservice \
	-- adapters::pgstorage::pgdb_test::pgstorage_test::test_search_by_key \
	--exact --show-output

# INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ('6411a28b-640a-43d9-b901-1c4b15d91568', 'frederick', 'long name', 'multiple names', 'names', 'name names');
# INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ('5a2579f7-83b9-4891-8dbc-e0024b5f3505', 'rick', 'short name', 'just one name', 'names', 'name names');
# INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) VALUES ('22cfc4fb-f9b6-4f6e-9158-9982347ad2a7', 'patrick', 'a saint', 'names', 'words', 'over words');
# SELECT KB_ID, KB_KEY, KIND, TAGS::TEXT AS TAGS FROM kbs WHERE KB_KEY LIKE '%rick%' LIMIT 10 OFFSET 0;
.PHONY: test-search
test-search : ## test search in the persistence layer.
	INTEGRATION_TEST=1 cargo test --package kbservice --bin kbservice \
	-- adapters::pgstorage::pgdb_test::pgstorage_test::test_search \
	--exact --show-output

.PHONY: test-save-kb
test-save-kb : ## test save key in the persistence layer.
	INTEGRATION_TEST=1 cargo test --package kbservice --bin kbservice \
	-- adapters::pgstorage::pgdb_test::pgstorage_test::test_save_kb \
	--exact --show-output

.PHONY: test-save-category
test-save-category : ## test save category in the persistence layer.
	INTEGRATION_TEST=1 cargo test --package kbservice --bin kbservice \
	-- adapters::pgstorage::pgdb_test::pgstorage_test::test_save_category \
	--exact --show-output --nocapture

# INSERT INTO categories (CATEGORY_NAME, CATEGORY_DESC) VALUES ('concepts', 'a principle or idea');
# INSERT INTO categories (CATEGORY_NAME, CATEGORY_DESC) VALUES ('software', 'any topic related to software development');
# INSERT INTO categories (CATEGORY_NAME, CATEGORY_DESC) VALUES ('personal', 'any personal data');
# SELECT CATEGORY_NAME, CATEGORY_DESC FROM categories ORDER BY CATEGORY_NAME LIMIT 5 OFFSET 0;
.PHONY: test-list-categories
test-list-categories : ## test listing categories without keywords
	INTEGRATION_TEST=1 cargo test --package kbservice --bin kbservice \
	-- adapters::pgstorage::pgdb_test::pgstorage_test::test_list_categories \
	--exact --show-output --nocapture

.PHONY: test-list-categories-with-keyword
test-list-categories-with-keyword : ## test listing categories without keywords
	INTEGRATION_TEST=1 cargo test --package kbservice --bin kbservice \
	-- adapters::pgstorage::pgdb_test::pgstorage_test::test_list_categories_with_keyword \
	--exact --show-output --nocapture

.PHONY: run-build-image
run-build-image: ## build container image
	${CONTAINERCMD} build \
	--build-arg appVersion=${VERSION} \
	--build-arg buildDate=${BUILD_DATE} \
	--build-arg commitHash=${COMMIT_HASH} \
    -f Dockerfile \
    -t kbservice:${VERSION} .
