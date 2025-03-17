# Key-Value Store

A simple in-memory key-value store implemented in Go. This project provides a basic REST API to store, retrieve, and delete keys, along with optional TTL (time-to-live) for keys.

## Features
- **In-Memory Storage:** Fast and lightweight key-value store using Go maps.
- **REST API Endpoints:**
  - `POST /store`: Add a new key-value pair (with optional TTL).
  - `GET /store/:key`: Retrieve the value associated with a key.
  - `DELETE /store/:key`: Delete a key from the store.
- **Time-to-Live (TTL):** Set a TTL for keys so they automatically expire.
- **Graceful Shutdown:** Handles `SIGTERM` and `SIGINT` signals for clean shutdown.
- **Logging:** Logs requests and responses to both stdout and a file (`logs.txt`).
- **Benchmark Testing:** Includes `testing.B` benchmarks to measure performance.

## Getting Started
### Prerequisites
- Go 1.19 or later
- Make (for convenient commands)

### Running the Server
1. Clone the repository:
   ```bash
   git clone https://github.com/NavroO/go-key-value-store.git
   cd go-key-value-store
   ```
2. Run the server:
   ```bash
   make run
   ```

The server will start on port `8080` by default.

### Testing
Run the unit tests:
```bash
make test
```

Run the benchmark tests:
```bash
make benchmark
```

### Logging
By default, logs are written to `logs.txt` and printed to stdout. You can review the log file for detailed information on API requests and responses.

## API Endpoints
| Method | Endpoint        | Description              |
|--------|------------------|--------------------------|
| POST   | `/store`         | Add a key-value pair     |
| GET    | `/store/:key`    | Retrieve a value by key  |
| DELETE | `/store/:key`    | Delete a key             |

### Example Requests
1. **Add a key-value pair:**
   ```bash
   curl -X POST http://localhost:8080/store \
        -H "Content-Type: application/json" \
        -d '{"key": "username", "value": "john_doe", "ttl": 60}'
   ```

2. **Retrieve a key:**
   ```bash
   curl http://localhost:8080/store/username
   ```

3. **Delete a key:**
   ```bash
   curl -X DELETE http://localhost:8080/store/username
   ```

## Contributing
Contributions are welcome! Feel free to open issues or submit pull requests.

## Acknowledgements
Thanks to everyone who contributed ideas and feedback for this project.

## Contact
- GitHub: [NavroO](https://github.com/NavroO)
