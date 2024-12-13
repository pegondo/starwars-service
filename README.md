# Starwars service

Starwars service is a simple Go service for interacting with [SWAPI - The Start Wars API](https://swapi.dev/).

## Features

Starwars service provides a REST API to interact with the [people](https://swapi.dev/documentation#people) and [planets](https://swapi.dev/documentation#planets) collections from [SWAPI](https://swapi.dev/).
Besides this basic interation, it also handles:

- **Pagination**: [SWAPI](https://swapi.dev/) only supports pagination with a page size of 10 elements, but Starwars service manages the pagination internally to serve pages of any size.
- **Search**: as [SWAPI](https://swapi.dev/), Starwars API supports search by name in both the [people](https://swapi.dev/documentation#people) and [planets](https://swapi.dev/documentation#planets) collections.
- **Sorting**: Starwars service extends the [SWAPI](https://swapi.dev/) functionality by sorting the results based on the `name` or `created` fields in `ascending` or `descending` order.

## Run the service

To run Starwars API, you have to first clone this respository:

```bash
git clone https://github.com/pegondo/starwars-service.git
cd starwars-service
```

Then, you have two options to run the service: as a go binary or as in a docker container.

### Run as a binary

To run the service as a binary:

1. Build the service:

```bash
go  build  -o  service  .
```

2. Run the binary:

```bash
./service
```

### Run in a docker container

To run the service in a docker container:

1. Build the docker image:

```bash
docker compose -f docker/docker-compose.yaml build
```

2. Start the container:

```bash
docker compose -f docker/docker-compose.yaml up
```

No matter the method you used, the service will be running in port `8080`.

## Endpoints

You can find the documentation for the endpoints in [this Swagger file](/docs/api/swagger/api.yaml).

## Testing

Starwars service has unit and integration tests.

### Run the unit tests

To run the unit tests, use:

```bash
go test ./...
```

### Run the integration tests

To run the integration tests:

1. [Run the service](#run-the-service).
2. Navigate to the `itests` folder and run the test command:

```bash
cd itests
go test ./...
```
