# Iran Access Domain Checker

A Go-based web application that monitors the accessibility of Iranian domains and provides a web interface to view their status.

## Features

- **Concurrent Domain Checking**: Checks multiple domains simultaneously for better performance
- **Automatic Monitoring**: Continuously monitors domains every 10 minutes
- **Web Interface**: Clean HTML/CSS/JS frontend to display domain status
- **RESTful API**: JSON API endpoint for domain status data
- **Thread-Safe Caching**: Concurrent-safe caching mechanism for domain results
- **Graceful Shutdown**: Proper server shutdown handling
- **Error Handling**: Comprehensive error handling throughout the application

## Project Structure

```
iranAccess/
├── server.go                     # Main application entry point
├── go.mod                        # Go module dependencies
├── go.sum                        # Go module checksums
├── domains.csv                   # CSV file containing domains to check
├── public/                       # Static web files
│   ├── index.html                # Main web interface
│   ├── script.js                 # Frontend JavaScript
│   └── style.css                 # Styling
└── internal/                     # Internal packages
    ├── models/
    │   └── domain.go              # Domain result data structures
    ├── cache/
    │   └── cache.go               # Thread-safe result caching
    ├── checker/
    │   └── checker.go             # Domain checking logic
    └── handlers/
        └── handlers.go            # HTTP request handlers
```

## API Endpoints

- `GET /check` - Returns the current status of all monitored domains
- `GET /api/check` - Alternative API endpoint for domain status

## Installation and Usage

1. **Prerequisites**: Go 1.24 or higher

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Prepare domains file**:
   Create a `domains.csv` file with the following format:
   ```csv
   domain
   google.com
   github.com
   stackoverflow.com
   ```

4. **Run the application**:
   ```bash
   go run server.go
   ```

5. **Access the web interface**:
   Open your browser and navigate to `http://localhost:3000`

## Configuration

- **Port**: Default port is 3000 (configurable in `server.go`)
- **Check Interval**: Domains are checked every 10 minutes (configurable in `checker.go`)
- **Request Timeout**: HTTP requests timeout after 5 seconds (configurable in `checker.go`)
- **CSV File Path**: Default is `./domains.csv` (configurable in `server.go`)

## Architecture Improvements

### Separation of Concerns
- **Models**: Data structures and types
- **Cache**: Thread-safe caching mechanism
- **Checker**: Domain checking business logic
- **Handlers**: HTTP request handling
- **Main**: Application setup and configuration

### Error Handling
- Proper error propagation and logging
- JSON error responses for API endpoints
- Graceful handling of CSV file errors

### Concurrency
- Thread-safe cache with RWMutex
- Concurrent domain checking with goroutines and WaitGroups
- Proper channel management

### Code Quality
- Consistent naming conventions
- Proper package organization
- Documentation and comments
- Separation of configuration constants

## Development

To extend this application:

1. **Add new endpoints**: Create handlers in `internal/handlers/`
2. **Modify checking logic**: Update `internal/checker/checker.go`
3. **Change data models**: Modify `internal/models/domain.go`
4. **Add middleware**: Include in `server.go` setup

## Dependencies

- [Fiber](https://github.com/gofiber/fiber) - Fast HTTP web framework
- Standard Go library for CSV parsing, HTTP client, and concurrency
