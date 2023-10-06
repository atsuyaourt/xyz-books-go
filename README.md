# XYZ Books

XYZ Books is a demo web application built in Go that offers the following features:

1. Database Schema Management, CRUD UI, and API.
2. A set of Go services that communicate with the web application via the API.

## Getting Started

To run the server, follow these steps:

1. Clone this repository.
2. Run the following command:

```console
docker compose up -d
```

By default, the application runs on port 3000.

## Endpoints

- `/`: Displays a list of available books with search functionality and pagination.
- `/{isbn13}`: Displays details for a book identified by its ISBN-13.
- `/api/v1`: The API endpoint (see below for more information).
- `/api/v1/docs/index.html`: Access the API documentation generated using [Swag](https://github.com/swaggo/swag).

## Database Schema

Schema can be found here: [db schema](db/migrations/000001_init_schema.up.sql). Codes for interacting with the database are generated using [sqlc](https://sqlc.dev/). Related queries and generated codes can be found in [db/query](db/query) and [db/sqlc](db/sqlc), respectively. Mocks for testing are located in [db/mocks](db/mocks) and are generated using [Mockery](https://vektra.github.io/mockery/latest/).

## JSON API

The JSON API is powered by [Gin](https://gin-gonic.com/). The [code](internal/api) includes CRUD handlers for book, author and publisher models.

## Front End

Front end is built with [Vite](https://v2.vitejs.dev/) [VueJS](https://vuejs.org/).

## Services

### ISBN

The [ISBN service](internal/service/isbn.go) performs the following tasks:

1. Call the books index endpoint.
2. Converts ISBN-10 to ISBN-13 and vice versa.
3. Updates missing ISBN-10 or ISBN-13 via the update endpoint.
4. Appends new ISBNs/EANs to a CSV file. _CSV file name is 'isbn.csv'_

## Environment Variables

See [.env.example](./.env.example)
