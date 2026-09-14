# CRM Backend

A simple customer relationship management REST API written in Go using only the standard library. It supports creating, listing, retrieving, updating, and deleting customers, including their contact details and contacted status. Data is stored in memory with four sample customers and resets when the server restarts.

## Installation

1. Install Go 1.22 or newer.
2. Clone or download this repository and open a terminal in the `CRM_Backend` directory.

No external dependencies or database setup are required.

## Launch

```sh
go run main.go
```

The API runs at `http://localhost:8080`. Press `Ctrl+C` to stop it.

## Usage

Use curl or an API client to call these endpoints. Keep the trailing slashes exactly as shown.

| Method | Endpoint | Action |
| --- | --- | --- |
| GET | `/customers` | List all customers |
| GET | `/customers/?id=1` | Get a customer by ID |
| POST | `/customers` | Create a customer |
| PUT | `/customers/` | Update a customer using the ID in the JSON body |
| DELETE | `/customer/?id=1` | Delete a customer by ID |

List customers:

```sh
curl http://localhost:8080/customers
```

Create a customer (all four fields are required; `id` is assigned automatically and `contacted` starts as `false`):

```sh
curl -X POST http://localhost:8080/customers \
  -H 'Content-Type: application/json' \
  -d '{"name":"Jane Doe","role":"Member","email":"jane@example.com","phone":"0123456789"}'
```

Update a customer by supplying the complete record:

```sh
curl -X PUT http://localhost:8080/customers/ \
  -H 'Content-Type: application/json' \
  -d '{"id":1,"name":"Bob","role":"Premium Member","email":"bob@example.com","phone":"0678304019","contacted":true}'
```

The current update validation requires `contacted` to be `true`.

Delete a customer (successful deletion returns HTTP 204 with no body):

```sh
curl -X DELETE 'http://localhost:8080/customer/?id=1'
```
