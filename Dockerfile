# ---------- Stage 1: Build ----------
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

# ---------- Stage 2: Runtime ----------
FROM debian:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main", "--port=8080", "--auth=true"]