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
│   │   └── user.go
│   ├── infrastructure/
│   │   └── db/
│   │       └── db.go
│   ├── middlewares/
│   │   ├── auth.go
│   │   └── cors.go
│   ├── models/
│   │   └── user.go
│   ├── service/
│   │   └── user.go
│   └── storage/
│       ├── errors.go
│       └── user.go
├── migrations/
│   ├── 20250221132128_add_users_table.sql
│   └── 20250221141157_add_users_table.sql
├── pkg/
├── Dockerfile
├── README.md
├── docker-compose.yaml
├── go.mod
└── go.sum

```

