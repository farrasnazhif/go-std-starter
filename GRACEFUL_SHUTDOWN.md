# Graceful Shutdown

The API implements graceful shutdown to ensure clean termination and proper resource cleanup.

## How It Works

When the server receives a shutdown signal (SIGINT or SIGTERM):

1. **Signal Detection**: The server detects the OS signal
2. **Stop Accepting Requests**: The server stops accepting new connections
3. **Wait for Active Requests**: Existing requests are allowed to complete
4. **Timeout**: If requests don't complete within 30 seconds, they're forcefully terminated
5. **Resource Cleanup**: Database connections and other resources are properly closed

## Signals Handled

- **SIGINT** - Triggered by `Ctrl+C` in terminal
- **SIGTERM** - Triggered by kill commands (`kill <pid>`, Docker container stop, etc.)

## Configuration

The graceful shutdown timeout is set to **30 seconds** in `cmd/api/shutdown.go`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
```

To customize the timeout, modify the timeout value:

```go
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // 60 seconds
```

## Example Usage

### Local Development

```bash
# Start the server
go run ./cmd/api/

# In another terminal, gracefully shutdown
kill -SIGTERM $(pgrep -f "go run ./cmd/api/")
# or simply press Ctrl+C
```

### Docker

```bash
# Stop container with graceful shutdown (SIGTERM)
docker stop container_name

# Stop container immediately (SIGKILL) - less graceful
docker kill container_name
```

## What Gets Closed

During graceful shutdown, the following resources are properly cleaned up:

1. **HTTP Server** - Stops accepting new connections
2. **Active Requests** - Allowed to complete within timeout
3. **Database Connections** - Connection pool is drained and closed
4. **Logger** - Buffered logs are flushed
5. **Email Service** - Any in-flight email operations complete

## Best Practices

1. **Always handle SIGTERM** - Cloud platforms and container orchestrators send SIGTERM, not SIGKILL
2. **Set appropriate timeout** - Balance between giving requests time to complete and fast shutdown
3. **Monitor logs** - Check for clean shutdown messages in logs
4. **Test gracefully** - Regularly test shutdown behavior during active requests

## Logging

Graceful shutdown logs key events:

```
"Received signal, shutting down gracefully" signal=SIGTERM
"Graceful shutdown initiated"
"Server shutdown completed successfully"
```

Or if an error occurs:

```
"Server shutdown error" error=<error message>
```

## Example Startup/Shutdown Log

```
INFO: Server starting addr=:8080
INFO: database connection pool established

[signal received]

INFO: Received signal, shutting down gracefully signal=SIGTERM
INFO: Graceful shutdown initiated
INFO: Server shutdown completed successfully
```

## Handling Force Shutdown

If the graceful shutdown times out or fails, the server will:

1. Log the error
2. Force close the HTTP server
3. Return the error
4. The process exits, triggering any OS-level cleanup
