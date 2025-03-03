# Astralis - Astronomical Events API

Astralis is a Go-based REST API and CLI application that provides information about upcoming astronomical events such as meteor showers, eclipses, and other celestial phenomena.

## Features

- REST API for querying astronomical events
- Filter events by date range and type
- CLI application with ASCII art visualization
- Hexagonal architecture for easy extension and maintenance
- Multiple data sources:
  - NASA DONKI API for solar events
  - Visible Planets API for planetary positions and visibility
- Comprehensive test suite

## Prerequisites

- Go 1.16 or later
- Environment variables for API configuration:
  - `NASA_API_KEY`: API key for NASA APIs (get one at https://api.nasa.gov/)
  - `PORT`: Port for the REST API server (default: 8080)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/alandavd/astralis.git
   cd astralis
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running Tests

The project includes a comprehensive test suite covering domain models, service layer, and API handlers.

Run all tests:

```bash
go test ./...
```

Run tests with coverage report:

```bash
go test -cover ./...
```

Generate detailed coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Running the API Server

1. Set the required environment variables:

   ```bash
   export NASA_API_KEY="your-nasa-api-key"
   export PORT="8080"
   ```

2. Start the server:

   ```bash
   go run cmd/api/main.go
   ```

   ```bash
   go run cmd/api/main.go -api_port=:8081
   ```

## Using the CLI Client

The CLI client provides a text-based interface with ASCII art visualization of astronomical events.

```bash
go run cmd/cli/main.go --api http://localhost:8080
```

## API Endpoints

- `GET /events`: Get all upcoming events

  - Query parameters:
    - `start`: Start date (RFC3339 format)
    - `end`: End date (RFC3339 format)

- `GET /events/{id}`: Get a specific event by ID

- `GET /events/type/{type}`: Get events by type
  - Supported types:
    - METEOR_SHOWER
    - ECLIPSE
    - CONJUNCTION
    - TRANSIT
    - OTHER

## Data Sources

### NASA DONKI API

- Provides solar events data (CMEs, solar flares)
- Free API key required (get one at https://api.nasa.gov/)
- Updates daily

### Visible Planets API

- Provides planetary visibility and position data
- No API key required
- Real-time calculations

## Architecture

The application follows hexagonal architecture principles:

- `internal/core/domain`: Domain entities and business logic
- `internal/core/ports`: Interface definitions
- `internal/core/service`: Business logic implementation
- `internal/adapters/primary`: Input adapters (REST API, CLI)
- `internal/adapters/secondary`: Output adapters (NASA API, Visible Planets API)

## Testing Strategy

The test suite follows these principles:

1. **Domain Tests**: Verify the core business logic and domain model behavior
2. **Service Tests**: Ensure the service layer correctly orchestrates repositories
3. **Handler Tests**: Validate HTTP endpoints and request/response handling
4. **Mock Objects**: Use mock implementations for external dependencies

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
