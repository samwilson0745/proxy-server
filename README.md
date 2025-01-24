# proxy-server

This is a simple caching proxy server written in Go. It forwards requests to an origin server and caches the responses for a specified time-to-live (TTL).

## Usage

### Running the Server

To run the server, use the following command:

```sh
go run cmd/main.go --port <port> --origin <origin-url>
```

- `--port`: The port on which the caching proxy server will run (default: 3000).
- `--origin`: The URL of the origin server to which requests will be forwarded.
- `--clear-cache`: Clears the cache and exits.

### Example

```sh
go run cmd/main.go --port 3000 --origin http://example.com
```

This will start the proxy server on port 3000 and forward requests to `http://example.com`.

### Clearing the Cache

To clear the cache, use the `--clear-cache` flag:

```sh
go run cmd/main.go --clear-cache
```

## Building the Server

To build the server, use the following command:

```sh
go build -o main cmd/main.go
```

This will create an executable named `main` in the current directory.

### Running After Building

To run the built server, use the following command:

```sh
./main --port <port> --origin <origin-url>
```

## Code Structure

- `cmd/main.go`: Entry point for the application.
- `server/server.go`: Contains the main server logic, including request handling and caching.
- `server/cache.go`: Contains the function to clear the cache.

## License

This project is licensed under the MIT License.