

````markdown
# 📚 Book Project

A Go-based RESTful API for managing books, built with clean architecture principles and supporting CRUD operations with optional JWT and Basic Authentication.

---

## 📖 Overview

The Book Project is a robust, modular API developed in Go for managing a collection of books. It provides endpoints for creating, reading, updating, and deleting books, with a clean architecture that separates concerns into domain, application, and infrastructure layers.

Authentication is optional, using **JWT** for protected endpoints and **Basic Authentication** for listing books and token generation. The `--auth=false` flag disables all authentication, making it ideal for development and testing.

---

## ✨ Features

### 🔧 CRUD Operations

- `GET /api/v1/books`: List all books.
- `POST /api/v1/books`: Create a new book.
- `GET /api/v1/books/{uuid}`: Retrieve a book by UUID.
- `PUT /api/v1/books/{uuid}`: Update a book by UUID.
- `DELETE /api/v1/books/{uuid}`: Delete a book by UUID.

### 🔐 Authentication

- Optional **JWT authentication** for `POST`, `GET /{uuid}`, `PUT`, and `DELETE` endpoints.
- **Basic Authentication** for:
  - `GET /api/v1/books`
  - `GET /api/v1/get-token`
- Disable authentication with `--auth=false`.

### 🧠 Architecture

- 🗂️ In-Memory Storage (easily extensible to databases)
- 🧱 Clean Architecture (handlers, services, repos, domain separated)
- 🖥️ CLI Support with Cobra

---

## 🛠️ Prerequisites

- **Go**: Version 1.18 or higher ([Install Go](https://go.dev/doc/install))
- **Git**: For cloning the repository ([Install Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git))
- **curl**: For testing API endpoints (optional)

---

## 🚀 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/biswasurmi/Book_Project.git
cd Book_Project
````

### 2. Install Dependencies

```bash
go mod tidy
```

**Dependencies include:**

* [`github.com/go-chi/chi/v5`](https://github.com/go-chi/chi)
* [`github.com/go-chi/jwtauth/v5`](https://github.com/go-chi/jwtauth)
* [`github.com/google/uuid`](https://github.com/google/uuid)
* [`github.com/spf13/cobra`](https://github.com/spf13/cobra)

### 3. Run the Server

* **With authentication:**

```bash
go run main.go startProject --port=8080 --auth=true
```

* **Without authentication:**

```bash
go run main.go startProject --port=8080 --auth=false
```

---

## 📡 API Usage

### 🔐 With Authentication (`--auth=true`)

#### 1. Get a JWT Token

```bash
curl -u admin:admin123 http://localhost:8080/api/v1/get-token
```

**Response:**

```json
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}
```

#### 2. Create a Book

```bash
curl -X POST http://localhost:8080/api/v1/books \
-H "Authorization: Bearer <your-jwt-token>" \
-H "Content-Type: application/json" \
-d '{"name":"Test Book","authorList":["Author One"],"publishDate":"2023-01-01","isbn":"1234567890"}'
```

#### 3. List Books

```bash
curl -u admin:admin123 http://localhost:8080/api/v1/books
```

#### 4. Get a Book

```bash
curl -H "Authorization: Bearer <your-jwt-token>" \
http://localhost:8080/api/v1/books/<uuid>
```

#### 5. Update a Book

```bash
curl -X PUT http://localhost:8080/api/v1/books/<uuid> \
-H "Authorization: Bearer <your-jwt-token>" \
-H "Content-Type: application/json" \
-d '{"name":"Updated Book","authorList":["Author Two"],"publishDate":"2024-01-01","isbn":"0987654321"}'
```

#### 6. Delete a Book

```bash
curl -X DELETE http://localhost:8080/api/v1/books/<uuid> \
-H "Authorization: Bearer <your-jwt-token>"
```

---

### 🔓 Without Authentication (`--auth=false`)

#### 1. Get a Token

```bash
curl http://localhost:8080/api/v1/get-token
```

#### 2. Create a Book

```bash
curl -X POST http://localhost:8080/api/v1/books \
-H "Content-Type: application/json" \
-d '{"name":"Test Book","authorList":["Author One"],"publishDate":"2023-01-01","isbn":"1234567890"}'
```

#### 3. List, Get, Update, Delete

Use the same commands as above, but omit the `Authorization` header and replace `<uuid>` with a valid UUID.

---

## 📂 Project Structure

```
Book_Project/
├── api/
│   ├── handler/         # HTTP handlers and routes
│   ├── middleware/      # Authentication middleware (JWT, Basic Auth)
├── cmd/                 # CLI commands using Cobra
├── domain/
│   ├── entity/          # Book entity definition
│   ├── repository/      # Repository interfaces
├── infrastructure/
│   ├── persistance/
│   │   ├── inmemory/    # In-memory repository implementation
├── service/             # Business logic for book operations
├── main.go              # Application entry point
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
├── LICENSE              # MIT License
├── README.md            # Project documentation
```
