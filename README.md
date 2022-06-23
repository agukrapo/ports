# Ports service

## Usage
A Postgres connection string must be provided in the `DATABASE_DSN` environment variable (see `.env.example`).

### CLI
`make build && ./bin/ports cli ports.json`

### REST server
`make build && ./bin/ports rest`

In other terminal

`curl -v -X PUT -F file=@ports.json localhost:8080/upload`

### gRPC server
`make build && ./bin/ports grpc-server`

In other terminal

`./bin/ports grpc-client localhost:8080 ports.json`

### Docker

`docker-compose run -v "$PWD:$PWD" ports cli $PWD/ports.json`

Note: you may need to replace `$PWD` with the current absolute path.

### Tests

`make test`

### Linter

`make lint`

Note: [golangci-lint](https://golangci-lint.run/usage/install/) must be installed.

## Missing features

* Proper unit test coverage for `database` and `service` packages.
* End to end service integration test using test containers.
