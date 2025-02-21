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
│   │   ├── auth.go
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
│   ├── 20250221132128_add_users_table.sql
│   └── 20250221141157_add_users_table.sql
├── pkg/
│   └── hash/
│       └── hash.go
├── Dockerfile
├── README.md
├── docker-compose.yaml
├── go.mod
└── go.sum
```

