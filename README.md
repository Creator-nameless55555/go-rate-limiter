<<<<<<< HEAD
# Go Rate Limiter

A simple, thread-safe HTTP rate limiter built with Go's standard library.

This project demonstrates how to limit the number of requests a client can make within a specified time window using custom HTTP middleware.

## Features

- Configurable request limit and time window
- Per-client rate limiting based on IP address
- Thread-safe request tracking
- HTTP middleware integration
- JSON error responses when the limit is exceeded
- `Retry-After` header indicating when the client can retry
- Automated tests for rate limiting and middleware behavior
- No external dependencies

## Requirements

- Go 1.26.1 or a compatible version supporting the project's `go.mod` configuration

## Getting Started

Clone the repository:

```bash
git clone https://github.com/YOUR_USERNAME/go-rate-limiter.git
cd go-rate-limiter
```

Replace `YOUR_USERNAME` with your GitHub username after publishing the repository.

Run the server:

```bash
go run .
```

The server starts at:

`http://localhost:8080`

## Usage

Send a request to the `/hello` endpoint:

```bash
curl -i http://localhost:8080/hello
```

Successful response:

```http
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

Hello! Your request was accepted.
```

The default configuration allows **5 requests per client every 10 seconds**.

Once the limit is exceeded, the server returns `429 Too Many Requests`.

Example error response:

```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Retry-After: 10
```

```json
{
  "error": "rate limit exceeded",
  "retry_after": 10
}
```

The `Retry-After` header specifies the number of seconds the client should wait before retrying. The JSON response includes the same retry delay.

## Configuration

The rate limit is configured in `main.go`:

```go
limiter := NewRateLimiter(5, 10*time.Second)
```

- `5` — maximum number of requests allowed per client during a window.
- `10*time.Second` — duration of the rate-limit window.

## Running Tests

Run all tests:

```bash
go test -v ./...
```

Run static analysis:

```bash
go vet ./...
```

Build the project:

```bash
go build ./...
```

The test suite covers request limits, client isolation, window expiration, middleware behavior, and rate-limit error responses.

## Project Structure

```text
go-rate-limiter/
├── main.go
├── limiter.go
├── middleware.go
├── limiter_test.go
├── go.mod
├── .gitignore
└── README.md
```

## Limitations

- Client identity is based on the remote IP address.
- Rate-limit state is stored in memory and resets when the server restarts.
- The limiter is intended as a learning project and is not a complete distributed rate-limiting solution.
- When deployed behind a reverse proxy, client IP detection requires additional configuration.

## Technologies

- Go
- `net/http`
- `sync.Mutex`
- `time`
- `encoding/json`
- Go testing and `httptest`

## License

This project is available under the MIT License if a corresponding `LICENSE` file is added.
=======
# go-rate-limiter
>>>>>>> 7fb92810351a7001b7bc2a696ef20bcdbf648bd5
