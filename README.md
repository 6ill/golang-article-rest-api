# Go Article REST API

A simple RESTful API for managing articles and authors, built with Go, Fiber, and PostgreSQL.

## Prerequisites

* Go (version 1.23 or higher)
* Docker and Docker Compose
* `make` command-line tool (optional, for running tests and migrations from host machine)
* `migrate` command-line tool (optional, for running migrations manually)

## Setup

1.  Clone the repository:

    ```bash
    git clone <repository_url>
    cd golang-article-rest-api
    ```

2.  Environment Variables:
    Copy the example environment file and update it with your desired configuration.

    ```bash
    cp .env.example .env
    ```
    Edit the [.env](.env) file. The `docker-compose.yaml` file uses variables defined here.

3.  Build and Run with Docker Compose:
    This will build the Go application image, set up the PostgreSQL database containers (one for the main app, one for tests), and start the services.

    ```bash
    docker-compose up --build -d
    ```

## Running Migrations

Database migrations are handled by the `migrate` tool and are applied automatically by Docker Compose on startup using the volumes defined in [docker-compose.yaml](docker-compose.yaml).

You can also run migrations manually using the provided Makefile targets:

*   Create a new migration:
    ```bash
    make db/migrations/new name=create_users_table # Replace create_users_table with your migration name
    ```

*   Apply pending migrations (main database):
    ```bash
    make db/migrations/up
    ```

*   Apply pending migrations (test database):
    ```bash
    make db/migrations/up/test
    ```

## Running Tests

Tests are located in the [internal/tests](internal/tests) directory. You can run them using the Makefile target:

```bash
make test