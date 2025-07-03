````markdown
# 📚 Book Project API

A robust Go-based RESTful API for managing books and users, built with Clean Architecture principles. Supports CRUD operations, JWT and Basic Authentication, unit testing, and Dockerized deployment.

---

## 🌟 Overview

The Book Project API is a secure, modular, and extensible Go API to manage books and user accounts. Built with Clean Architecture, it separates core logic into domain, service, and handler layers, and supports:

- JWT authentication for protected routes
- Basic Auth for token generation and book listing
- Option to disable auth in development (`--auth=false`)
- Cobra CLI for server management
- Unit testing and Docker support

---

## ✨ Feature Summary

| Resource | Method | Endpoint | Auth Required (`--auth=true`) | Auth-Free Mode (`--auth=false`) |
|----------|--------|----------|-------------------------------|---------------------------------|
| 📘 Books | GET    | `/api/v1/books`            | ✅ Basic Auth required           | ✅ No Auth                      |
| 📘 Books | POST   | `/api/v1/books`            | ✅ Bearer Token (JWT)            | ✅ No Auth                      |
| 📘 Books | GET    | `/api/v1/books/{uuid}`     | ✅ Bearer Token (JWT)            | ✅ No Auth                      |
| 📘 Books | PUT    | `/api/v1/books/{uuid}`     | ✅ Bearer Token (JWT)            | ✅ No Auth                      |
| 📘 Books | DELETE | `/api/v1/books/{uuid}`     | ✅ Bearer Token (JWT)            | ✅ No Auth                      |
| 👤 Users | POST   | `/api/v1/register`         | ❌ Open to all                   | ❌ Open to all                  |
| 👤 Users | POST   | `/api/v1/login`            | ❌ Open to all                   | ❌ Open to all                  |
| 👤 Users | GET    | `/api/v1/users/{id}`       | ✅ Bearer Token (JWT)            | ✅ No Auth                      |
| 👤 Users | GET    | `/api/v1/users/me`         | ✅ Bearer Token (JWT)            | ✅ No Auth                      |
| 👤 Users | PUT    | `/api/v1/users/{id}`       | ✅ Bearer Token (JWT)            | ✅ No Auth                      |
| 👤 Users | DELETE | `/api/v1/users/{id}`       | ✅ Bearer Token (JWT)            | ✅ No Auth                      |
| 🔐 Auth  | GET    | `/api/v1/get-token`        | ✅ Basic Auth required           | ✅ No Auth                      |

---

## 🔐 Authentication Modes

- **Auth Enabled (`--auth=true`)**
  - `GET /books` and `GET /get-token`: Require Basic Auth (`urmi`).
  - Other endpoints require a JWT token in `Authorization` header.

- **Auth Disabled (`--auth=false`)**
  - All routes are open. You can test without any auth headers.

---

## 🧠 Data Models

### 📘 Book

```json
{
  "uuid": "123e4567-e89b-12d3-a456-426614174001",
  "name": "Learn API",
  "authorList": ["author1", "author2"],
  "publishDate": "2022-01-02",
  "isbn": "0999-0555-5914"
}
````

### 👤 User

```json
{
  "firstName": "urmi",
  "lastName": "admin",
  "userName": "urmi",
  "password": "password123"
}
```

---

## 🚀 Usage Examples

### ✅ Register User

```bash
curl -X POST http://localhost:8080/api/v1/register \
-H "Content-Type: application/json" \
-d '{"firstName":"urmi","lastName":"admin","userName":"urmi","password":"password123"}'
```

---

### ✅ Login for JWT

```bash
curl -X POST http://localhost:8080/api/v1/login \
-H "Content-Type: application/json" \
-d '{"email":"urmi@example.com","password":"password123"}'
```

---

### ✅ Get Token (Basic Auth)

```bash
curl -u urmi:password123 http://localhost:8080/api/v1/get-token
```

---

### ✅ List Books (Basic Auth or Open)

```bash
curl -u urmi:password123 http://localhost:8080/api/v1/books
```

Or without auth if `--auth=false`:

```bash
curl http://localhost:8080/api/v1/books
```

---

### ✅ Create Book

```bash
curl -X POST http://localhost:8080/api/v1/books \
-H "Authorization: Bearer <your-jwt-token>" \
-H "Content-Type: application/json" \
-d '{"name":"Learn API","authorList":["author1","author2"],"publishDate":"2022-01-02","isbn":"0999-0555-5914"}'
```

---

## 🐳 Docker Setup

### 🧱 Build Image

```bash
docker build -t book-project-api:latest .
```

### 🚀 Run Container

```bash
docker run -p 8080:8080 book-project-api:latest
```

Without auth:

```bash
docker run -p 8080:8080 book-project-api:latest --auth=false
```

---

## 📁 Project Structure

```
Book_Project/
├── api/
│   ├── handler/         # Route handlers
│   ├── middleware/      # JWT & Basic auth
├── cmd/                 # Cobra CLI commands
├── domain/
│   ├── entity/          # Book & User models
│   ├── repository/      # Interfaces
├── infrastructure/
│   └── persistance/
│       └── inmemory/    # In-memory repo
├── service/             # Business logic
├── test_file/           # Unit tests
├── main.go              # Entry point
├── Dockerfile           # Multi-stage Dockerfile
├── go.mod / go.sum      # Go modules
```

---

## 🧪 Run Tests

```bash
go test -v ./test_file
```

Covers:

* 📘 Book endpoints
* 👤 User endpoints
* 🔐 Auth logic

---

## 🐳 Dockerfile

```dockerfile
# Build stage
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Runtime stage
FROM debian:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main", "--port=8080", "--auth=true"]
```

---

## 📚 Resources

* [Go JWT (jwtauth)](https://github.com/go-chi/jwtauth)
* [Cobra CLI](https://github.com/spf13/cobra)
* [Go HTTP Testing](https://go.dev/doc/tutorial/add-a-test)
* [Learn REST APIs](https://developer.mozilla.org/en-US/docs/Web/HTTP)

---

## 🙌 Contributing

1. Fork this repo
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Commit changes: `git commit -m "Add my feature"`
4. Push: `git push origin feature/my-feature`
5. Open a pull request

---

Built with ❤️ by **Urmi Biswas** – Happy coding! 🚀

```
