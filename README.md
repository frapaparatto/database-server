# database-server

An in-memory key-value store served over HTTP, written in Go with the standard library only.

## Requirements

Go 1.27.1 or later.

## Run

```sh
go run .
```

The server listens on `http://localhost:4000`.

## API

Both endpoints accept only `GET` (and `HEAD`, which a `GET` route also matches).

Responses from the handlers are JSON. Errors have the shape `{"error": "<message>"}`. Other methods get `405 Method Not Allowed` and unknown paths get `404 Not Found`; both come from the mux as plain text.

In zsh, quote the URLs in `curl` commands, because `?` is a glob character.

### `/set`

Stores one or more key-value pairs passed in the query string.

```sh
curl "localhost:4000/set?a=1"
# 200 {"a":"1"}

curl "localhost:4000/set?x=1&y=2"
# 200 {"x":"1","y":"2"}
```

The response contains the pairs just stored. Setting an existing key overwrites its value.

A request is stored entirely or not at all: if any pair is invalid, the server returns `400` and stores nothing.

| Request | Status | Body |
|---|---|---|
| `/set` | 400 | `{"error":"no parameters provided"}` |
| `/set?=1` | 400 | `{"error":"empty key is not allowed"}` |
| `/set?a=` | 400 | `{"error":"empty value for key 'a'"}` |
| `/set?a=1&a=2` | 400 | `{"error":"key 'a' must appear exactly once"}` |

### `/get`

Returns the value stored for `key`.

```sh
curl "localhost:4000/get?key=a"
# 200 {"a":"1"}
```

| Request | Status | Body |
|---|---|---|
| `/get?key=missing` | 404 | `{"error":"key 'missing' not found"}` |
| `/get` or `/get?key=` | 400 | `{"error":"missing or empty 'key' parameter"}` |
| `/get?key=a&foo=1` | 400 | `{"error":"only the 'key' parameter is allowed"}` |
| `/get?key=a&key=b` | 400 | `{"error":"'key' must appear exactly once"}` |

## Storage

Data lives in memory only and is lost when the server stops.

The store is a map protected by a mutex, so concurrent requests are safe. All pairs of one `/set` request are written under a single lock, so a reader never sees half of a request.

## Tests

```sh
go test -race ./...
```

- `TestGet` and `TestSet` check every case in the tables above: status, body, and for `/set` the content of the store after the request.
- `TestConcurrentAccess` runs reads and writes on the same store at the same time. With `-race`, it fails if any access to the store is not synchronized.

The `experiment` script sends requests to a running server and logs each response to `experiments.log`, for manual checks:

```sh
./experiment [label]
```

## Project layout

- `main.go`: creates the store, registers the routes, starts the server
- `handlers.go`: the `/get` and `/set` handlers
- `storage.go`: the store and its mutex
- `response.go`: helpers that write JSON and error responses
- `server_test.go`: tests

## Possible next steps

- Persist writes to a file, appending each write to a log, and rebuild the store from it at startup.
- Make it more efficient: the log grows with every write, so compact it periodically.

## AI use

I used AI to review my code, rewrite error messages and comments in clearer English, write commit messages, and draft this README from a structure I defined.
