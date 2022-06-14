# Ports service

## Usage

`make build && ./bin/ports ports.json`

Note: a Postgres connection string must be provided in the `DATABASE_DSN` environment variable (see `.env.example`).

### Docker

`docker-compose run -v "$PWD:$PWD" ports $PWD/ports.json`

Note: you may need to replace `$PWD` with the current absolute path.

### Tests

`make test`

### Linter

`make lint`

Note: [golangci-lint](https://golangci-lint.run/usage/install/) must be installed.

## Missing features

* Proper unit test coverage for `database` and `service` packages.
* End to end service integration test using test containers.
