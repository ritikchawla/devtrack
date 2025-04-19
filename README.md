# DevTrack API

This is the backend API for the DevTrack application, a tool designed to help developers track their projects, tasks, and progress.

## Prerequisites

Before you begin, ensure you have the following installed:

*   [Go](https://golang.org/doc/install) (version 1.24 or later recommended)
*   [Docker](https://docs.docker.com/get-docker/)
*   [Docker Compose](https://docs.docker.com/compose/install/)
*   [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) (for generating Go code from SQL)
*   [golang-migrate/migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) CLI (for manual migration management if needed, though Docker handles it on startup)
*   [Make](https://www.gnu.org/software/make/) (Optional, for potentially adding helper commands later)

## Setup

1.  **Clone the repository:**
    ```bash
    git clone <your-repository-url>
    cd devtrack
    ```

2.  **Create Environment File:**
    Copy the example environment file and configure your database settings:
    ```bash
    cp .env.example .env
    ```
    Edit the `.env` file with your desired database credentials (especially `POSTGRES_PASSWORD`). The defaults match the `docker-compose.yml` setup.

    ```dotenv
    # .env file
    GIN_MODE=debug
    SERVER_ADDRESS=0.0.0.0:8080

    # Database Configuration
    DB_DRIVER=postgres
    POSTGRES_USER=root
    POSTGRES_PASSWORD=password # Change this in production!
    POSTGRES_DB=devtrack
    DB_HOST=devtrack_postgres # Service name in docker-compose
    DB_PORT=5432
    SSL_MODE=disable

    # Combine into DB_SOURCE for migrate tool
    DB_SOURCE=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB}?sslmode=${SSL_MODE}
    ```

## Running the Application

This project uses Docker Compose to manage the API service and the PostgreSQL database.

1.  **Build and Run:**
    ```bash
    docker compose up --build
    ```
    This command will:
    *   Build the Go application binary inside a Docker container.
    *   Build the final API image, including the `migrate` tool.
    *   Start the PostgreSQL database container.
    *   Start the API container.
    *   The API container's entrypoint script will automatically run database migrations before starting the API server.

2.  **Accessing the API:**
    The API will be available at `http://localhost:8080`.

3.  **Stopping the Application:**
    ```bash
    docker compose down
    ```

## Database Migrations

Migrations are handled by `golang-migrate/migrate`.

*   **Migration Files:** Located in the `db/migration` directory. Use the format `<timestamp>_<name>.up.sql` and `<timestamp>_<name>.down.sql`.
*   **Automatic Application:** Migrations are automatically applied when the `api` service starts via the `entrypoint.sh` script.
*   **Manual Management (Optional):** If you need to run migrations manually (e.g., during development without Docker), ensure the `migrate` CLI is installed and use commands like:
    ```bash
    # Ensure DB_SOURCE is set in your environment or pass it directly
    migrate -path db/migration -database "$DB_SOURCE" up
    migrate -path db/migration -database "$DB_SOURCE" down 1
    ```

## SQLC (SQL Compiler)

We use `sqlc` to generate type-safe Go code for database interactions.

*   **Schema:** Defined in `db/migration/*.up.sql` files.
*   **Queries:** Written in `db/query/*.sql` files.
*   **Configuration:** `sqlc.yaml` defines where to find schema/queries and where to output Go code (`db/sqlc/`).
*   **Generating Code:** Run `sqlc generate` in the project root after modifying schema or queries.

## API Endpoints

*   **`POST /users`**: Create a new user.
    *   **Request Body:**
        ```json
        {
          "username": "string",
          "email": "string (valid email format)",
          "password": "string",
          "full_name": "string"
        }
        ```
    *   **Response (Success: 200 OK):**
        ```json
        {
          "id": "integer",
          "username": "string",
          "email": "string",
          "full_name": "string",
          "created_at": "timestamp"
        }
        ```
*   **`GET /health`**: Basic health check endpoint.
