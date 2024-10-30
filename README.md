# Go In-Memory Cache Service

A lightweight, RESTful in-memory caching service built with Go. This service supports basic cache operations with time-to-live (TTL) management, making it ideal for temporary data storage where persistence is not required.

## API Endpoints

Below is an overview of the available endpoints:

![Endpoints Diagram](https://github.com/user-attachments/assets/934a357c-b02e-4c1f-9d36-f32ca4a11a5d)

### 1. **POST** `/cache`

Creates a new cache entry.

**Payload:**
- `key`: (string) The unique identifier for the cache entry.
- `ttl`: (string) Time-to-live for the cache entry in a format like "10min" for 10 minutes.
- `value`: (string) The data to store in the cache.

### 2. **POST** `/cache/get-or-set`

Retrieves an existing cache entry or creates a new one if it does not exist.

**Payload:**
- `key`: (string) The unique identifier for the cache entry.
- `ttl`: (string) Time-to-live for the cache entry, e.g., "10min".
- `value`: (string) The data to store if the cache entry does not already exist.

If a cache entry with the specified key already exists, the existing value will be returned. Otherwise, a new entry is created.

### 3. **GET** `/cache/:key`

Fetches the cached data associated with the specified `key`. Returns the data if found, otherwise responds with an appropriate error.

### 4. **DELETE** `/cache/:key`

Deletes the cache entry associated with the specified `key`.

## Future Enhancements

- **Benchmarking**: Implement benchmarking to measure and optimize performance under different loads and configurations.
