# Ports Service
This service allow users to import ports and create/update it

## Prerequisites

- Docker (for Docker-based setup)
- docker-compose
- Go (for local setup)
- Make

## Running the Server

### Using Docker

To run the server using Docker, follow these steps:

```bash
make create-network
make build
make up-dependencies
make migrate-up
make run-server
```
This will:
  - Create a network
  - Build the Docker image.
  - Set up the database.
  - Run migrations.
  - Start the server on port 8080.

### Locally
To run the server locally using the Go CLI, follow these steps:
```bash
make create-network
make up-dependencies
make migrate-up
make run-server-locally
```
This will:
  - Create a network
  - Set up the database.
  - Run migrations.
  - Start the server locally on port 8080.

## Running Imports

### Using Docker

To run imports using Docker, follow these steps:
```bash
make create-network
make build
make up-dependencies
make migrate-up
make run-import FILE=input/ports.json
```
*Note*: Ensure that the file you want to import is placed inside the /input directory.

This will:
  - Create a network
  - Build the Docker image.
  - Set up the database.
  - Run migrations.
  - Execute the import process using the specified JSON file.

### Locally

To run imports locally using the Go CLI, follow these steps:

```bash
make create-network
make up-dependencies
make migrate-up
make run-import-locally FILE=input/ports.json
```

*Note*: Ensure that the file you want to import is placed inside the /input directory.

This will:
  - Create network
  - Set up the database.
  - Run migrations.
  - Execute the import process locally using the specified JSON file.


### Utilities commands

```bash
make integration-tests
make unit-tests
make tests
make lint
```
Additional commands check the Makefile

## API

### Endpoints

### Create
`POST`: `localhost:8080/ports`
`payload`: 
```json
{
    "name": "NAME",
    "city": "Changshu",
    "province": "Jiangsu",
    "country": "China",
    "alias": [
      "Zhangjiagang",
      "Suzhou",
      "Taicang"
    ],
    "regions": [],
    "coordinates": [
      120.752503,
      31.653686
    ],
    "timezone": "Asia/Shanghai",
    "unlocs": [
      "CNCGA"
    ],
    "code": "57076"
  }
```
`reponses`: `200 OK` or `500 internal server error` 
*Note*: This API return 200 because this endpoint save or update if this port already exists

### Curl
```
curl --request POST \
  --url http://localhost:8080/ports \
  --header 'content-type: application/json' \
  --data '{
    "name": "saas",
    "city": "Changshu",
    "province": "Jiangsu",
    "country": "China",
    "alias": [
      "Zhangjiagang",
      "Suzhou",
      "Taicang"
    ],
    "regions": [],
    "coordinates": [
      120.752503,
      31.653686
    ],
    "timezone": "Asia/Shanghai",
    "unlocs": [
      "CNCGa"
    ],
    "code": "57076"
  }'
```

### GET by ID
`GET`: `localhost:8080/ports/{port_id}`

`response`:
```json
{
  "id": "CNCGA",
  "name": "saas",
  "city": "Changshu",
  "country": "China",
  "alias": [
    "Zhangjiagang",
    "Suzhou",
    "Taicang"
  ],
  "regions": [],
  "coordinates": [
    120.752503,
    31.653686
  ],
  "province": "Jiangsu",
  "timezone": "Asia/Shanghai",
  "unlocs": [
    "CNCGA"
  ],
  "code": "57076"
}
```

`http codes`: `200 OK`, `500 internal server error` or `404 not found`

*Note*: The API use the unloc code as ID for the elements

### Curl
``` 
curl --request GET \
--url http://localhost:8080/ports/CNCGa
```

## Architecture Overview
This project follows a hybrid approach, combining elements of **Clean Architecture** and **Hexagonal Architecture** to achieve a highly modular, maintainable, and scalable design. By structuring the code into well-defined layers—**Domain, Application, and Infrastructure**—we ensure a clear separation of concerns and strict dependency inversion.

### Key Architectural Principles:
- **Hexagonal Architecture (Ports & Adapters)**: The application core remains independent of external frameworks, using well-defined interfaces (ports) to interact with the outside world. This makes the system highly testable and adaptable to different infrastructures.
- **Clean Architecture Principles**: The domain and business logic are kept at the center, ensuring they remain free from dependencies on frameworks, databases, or delivery mechanisms.
- **Dependency Inversion**: High-level modules (business logic) do not depend on low-level modules (infrastructure); instead, both depend on abstractions.
- **Scalability & Maintainability**: The modular structure enables easier feature expansion, testing, and refactoring without impacting core business rules.

By leveraging these principles, the project achieves a strong balance between flexibility and maintainability, providing a solid foundation for evolving business needs and technological advancements.

## Project Structure
``` 
├── build
│   ├── docker-compose-app.yml
│   └── docker-compose-db.yml
├── cmd
│   ├── cli
│   │   ├── import.go
│   │   ├── root.go
│   │   └── server.go
│   └── main.go
├── config
│   └── config.go
├── Dockerfile
├── go.mod
├── go.sum
├── input
│   └── ports.json
├── internal
│   ├── core
│   │   ├── application
│   │   │   └── service.go
│   │   │   └── service_test.go
│   │   │   └── service_integration_test.go
│   │   └── domain
│   │       ├── ports.go
│   │       ├── errors.go
│   │       └── domain.go
│   └── infra
│       ├── adapters
│       │   ├── parser
│       │   │   ├── jsonparser.go
│       │   │   └── jsonparser_test.go
│       │   └── repository
│       │       ├── postgres.go
│       │       └── postgres_integration_test.go
│       └── server
│           └── http
│               └── handler
│                   ├── handler.go
│                   └── handler_integration_test.go
├── Makefile
├── migrations
│   ├── 000001_create_ports_table.down.sql
│   └── 000001_create_ports_table.up.sql
├── mocks
│   ├── parser_port.go
│   ├── repository_port.go
│   └── service_port.go
├── pkg
│   ├── database
│   │   └── postgres.go
│   └── graceful
│       └── shutdown.go
├── README.md
└── tests
    └── suite
        └── postgrescontainer.go
```

### `build/`
Contains Docker Compose files for setting up the application and the database.
- **`docker-compose-app.yml`**: Manages the application container.
- **`docker-compose-db.yml`**: Manages the database container.

### `cmd/`
The entry points of the application.
- **`cli/`**: Contains CLI-related commands.
    - `import.go`: Handles data import functionality.
    - `root.go`: Defines the root command for the CLI.
    - `server.go`: Manages server-related CLI commands.
- **`main.go`**: The main entry point of the application.

### `config/`
Manages application configuration.
- **`config.go`**: Handles loading and parsing configuration settings.

### `Dockerfile`
Defines the instructions for building the application's Docker image.

### `go.mod` & `go.sum`
Go module files managing dependencies and package versions.

### `input/`
Stores input data.
- **`ports.json`**: Contains port-related input data.

### `internal/`
The core business logic and infrastructure of the application.

#### `core/`
The business logic layer.
- **`application/`**: Contains service and application-specific ports.
    - `service.go`: Implements application services.
    - `service_test.go`: Unit tests for services.
    - `service_integration_test.go`: Integration tests services.
- **`domain/`**: Contains domain entities and business rules.
    - `ports.go`: Defines application ports (interfaces).
    - `errors.go`: Defines domain-specific errors.
    - `domain.go`: Represents the domain entity for ports.

#### `infra/`
Infrastructure-related implementations.
- **`adapters/`**: Connects external systems to the application.
    - **`parser/`**: Handles JSON parsing.
        - `jsonparser.go`: Implements JSON parsing logic.
        - `jsonparser_test.go`: Unit tests for JSON parsing.
    - **`repository/`**: Manages database interactions.
        - `postgres.go`: Implements PostgreSQL repository logic.
        - `postgres_integration_test.go`: Integration tests for PostgreSQL repository.
- **`server/`**: Handles HTTP server-related functionality.
    - **`http/handler/`**: Defines HTTP handlers.
        - `handler.go`: Implements request handling logic.
        - `handler_integration_test.go`: Integration tests for handlers.

### `Makefile`
Contains automation scripts for building, testing, and running the application.

### `migrations/`
Contains database migration scripts.
- **`000001_create_ports_table.up.sql`**: Creates the `ports` table.
- **`000001_create_ports_table.down.sql`**: Rolls back the `ports` table creation.

### `mocks/`
Stores mock implementations for unit testing.
- **`parser_port.go`**: Mock for the parser.
- **`repository_port.go`**: Mock for the repository.
- **`service_port.go`**: Mock for the service.

### `pkg/`
Reusable utilities and libraries.
- **`database/`**: Database-related utilities.
    - `postgres.go`: Provides PostgreSQL connection management.
- **`graceful/`**: Handles graceful shutdown of the application.
    - `shutdown.go`: Implements clean shutdown logic.

### `README.md`
Documentation for the project, including setup and usage instructions.

### `tests/`
Contains test-related utilities.
- **`suite/`**: Defines test suites.
    - `postgrescontainer.go`: Manages PostgreSQL container for integration testing.
