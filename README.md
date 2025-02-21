# file_storage

## Project Structure

```
├── cmd/
│   └── app/
│       └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── controller/
│   │   ├── errors.go
│   │   ├── interfaces.go
│   │   └── user.go
│   ├── infrastructure/
│   │   └── db/
│   │       └── db.go
│   ├── middlewares/
│   │   └── cors.go
│   ├── models/
│   │   └── user.go
│   ├── repository/
│   │   ├── errors.go
│   │   └── user.go
│   └── service/
│       ├── interfaces.go
│       └── user.go
├── migrations/
├── pkg/
├── Dockerfile
├── README.md
├── docker-compose.yaml
├── go.mod
└── go.sum
```

