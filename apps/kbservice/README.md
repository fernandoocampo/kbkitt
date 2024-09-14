# kb service

This is the core app for the kb system. This service will be in charge of storing kb entries and provide API to receive requests from different devices.

## Technologies

### Tokio
[Tokio](https://tokio.rs) is an asynchronous runtime for the Rust programming language. It provides the building blocks needed for writing network applications.

### Warp
[warp](https://docs.rs/warp/latest/warp/) is a super-easy, composable, web server framework for warp speeds.

## How to test it?

I added a Makefile, but feel free to run them as they are using cargo. I tried to add unit tests as much as possible to learn how to use them, I found many issues trying to test Warp due to lack of documentation, but I managed to do it.

Rust conventions suggest adding the unit tests in the same file where you have your code logic, but we're used to adding different files for that purpose, so I followed the Go convention.

```sh
make test
```

or

```sh
cargo test
```

## How to run it?

```sh
make run
```

or

```sh
RUST_LOG=debug cargo run
```

you will see something like this

```sh
➜  make run
RUST_LOG=debug LOG_SYSTEM=log4rs cargo run
    Finished dev [unoptimized + debuginfo] target(s) in 0.22s
     Running `target/debug/kbcore`
⏱️	Starting kbcore api application...
🪵	Initializing logger...
LOG_SYSTEM: log4rs
2023-09-30T19:58:16.219960+02:00 INFO kbcore::application::app - 🪵	Using log4rs
```

once you finished just hit `ctrl + c`

* another possible values for RUST_LOG

```log
error
warn
info
debug
trace
```

## How to call the API?

* Create a kb endpoint

```sh
curl -H "Content-Type: application/json" \
--data '{"name":"Esme"}' \
-X POST http://localhost:3030/people

{"id":"f1601fc5-f0c9-4950-8017-e094b284cad9"}
```

* Calling get people endpoint

```sh
curl -X GET http://localhost:3030/people
[{"id":"27c6bdd9-67d6-4503-884d-d75aba152f44","first_name":"Esme","last_name":"Esme"},{"id":"d49aed14-b5b0-4e49-972f-f823817ed93d","first_name":"Fernando","last_name":"Fernando"}]
```

with params

limit: The index of the last item which has to be returned
offset: The index of the first item which has to be returned

```sh
curl -X GET 'http://localhost:3030/people?limit=10&offset=0'

[{"id":"27c6bdd9-67d6-4503-884d-d75aba152f44","first_name":"Esme","last_name":"Esme"},{"id":"d49aed14-b5b0-4e49-972f-f823817ed93d","first_name":"Fernando","last_name":"Fernando"}]
```

* Calling get a person endpoint

```sh
curl -X GET http://localhost:3030/people/27c6bdd9-67d6-4503-884d-d75aba152f44

{"id":"27c6bdd9-67d6-4503-884d-d75aba152f44","first_name":"Esme","last_name":"Esme"}
```

```sh
curl -X GET http://localhost:3030/people/200
Person not found
```

* Update a person endpoint

```sh
curl -H "Content-Type: application/json" \
--data '{"id":"27c6bdd9-67d6-4503-884d-d75aba152f44", "first_name":"Esme", "last_name":"Emse"}' \
-X PUT http://localhost:3030/people

{"id":"27c6bdd9-67d6-4503-884d-d75aba152f44","first_name":"Esme","last_name":"Emse"}
```

* Delete a person endpoint

```sh
curl -H "Content-Type: application/json" \
-X DELETE http://localhost:3030/people/d3bc8246-53da-4275-b833-5feb4489741d

Person d3bc8246-53da-4275-b833-5feb4489741d deleted
```

## Migration

Project is using [sqlx-cli](https://docs.rs/crate/sqlx-cli/latest), so let's install it first.

```sh
cargo install sqlx-cli
```

* add migration for kbs table

```sh
sqlx migrate add -r create_kb_table

Creating migrations/20240503161844_create_kb_table.up.sql
Creating migrations/20240503161844_create_kb_table.down.sql
```

migration files were added in the `migrations` directory.

```sh
migrations/20230917172957_people_table.up.sql
migrations/20230917172957_people_table.down.sql
```

* run migrations

```sh
sqlx migrate run --database-url postgresql://localhost:5432/kbs
```

* revert migrations

Each revert will trigger the latest migration and try to run the `*.down.sql` script.

```sh
sqlx migrate revert --database-url "postgresql://localhost:5432/kbs"
```

## How to check database

* get into the database
```sh
psql -U kb_user -h localhost -p 5432
```

* list tables
```sh
kbs=# \dt
        List of relations
 Schema |  Name  | Type  | Owner
--------+--------+-------+-------
 public | kbs | table | kb_owner
(1 rows)
```

## TSVECTOR in Postgresql

project is using [TSVECTOR](https://www.postgresql.org/docs/current/datatype-textsearch.html) data type to make possible to search by keywords.

* How to insert a record

```sql
INSERT INTO kbs (KB_ID, KB_KEY, KB_VALUE, NOTES, KIND, TAGS) 
 	VALUES ('e9efcab6-adc0-4987-af18-49ca6da35f87', 'green', 'other color', 'to remember other color', 'concepts', 'color green paint concepts');
```

* How to query

```sql
SELECT * FROM kbs WHERE TAGS @@ to_tsquery('green');
```