FROM golang:1.23.5

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 8081

CMD ["go", "run", "cmd/app/main.go"]
