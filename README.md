# Ports service

## Usage
A Postgres connection string must be provided in the `DATABASE_DSN` environment variable (see `.env.example`).

### CLI
`make build && ./bin/ports cli ports.json`

### REST server
`make build && ./bin/ports rest`

In other terminal

`curl -v -X PUT -F file=@ports.json localhost:8080/upload `

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
