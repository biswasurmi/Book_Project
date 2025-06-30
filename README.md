

---


````md
# 📚 Book Management REST API (Go + Chi + JWT + Basic Auth)

A simple REST API for managing books with basic CRUD operations, written in Go using the Chi router.  
Supports **Basic Authentication** and **JWT-based authorization**. Data is stored in-memory (non-persistent).

---

## 🚀 Features

- Create, read, update, and delete books
- Basic Auth for token retrieval and listing
- JWT Auth for creating/updating/deleting
- In-memory storage (for simplicity)

---

## 📘 Book Model

```json
{
  "uuid": "auto-generated string",
  "name": "string",
  "authorList": ["string"],
  "publishDate": "YYYY-MM-DD",
  "isbn": "string"
}
````

---

## ⚙️ Getting Started

### ✅ Prerequisites

* Go 1.22 or higher
* Git (optional, for cloning)

### ▶️ Run Server

```bash
go run main.go
```

#### Command-line flags:

| Flag    | Default | Description                      |
| ------- | ------- | -------------------------------- |
| `-auth` | true    | Enable or disable authentication |
| `-port` | 8080    | Port to run the server on        |

**Example (disable authentication):**

```bash
go run main.go -auth=false
```

---

## 🔌 API Endpoints

| Method | Endpoint               | Description         | Auth Required     |
| ------ | ---------------------- | ------------------- | ----------------- |
| GET    | `/api/v1/get-token`    | Get JWT token       | Basic Auth        |
| GET    | `/api/v1/books`        | List all books      | Basic or JWT Auth |
| POST   | `/api/v1/books`        | Create a new book   | JWT               |
| GET    | `/api/v1/books/{uuid}` | Get book by UUID    | JWT               |
| PUT    | `/api/v1/books/{uuid}` | Update book by UUID | JWT               |
| DELETE | `/api/v1/books/{uuid}` | Delete book by UUID | JWT               |

---

## 💡 Usage Examples

### 🔐 Get JWT Token (via Basic Auth)

```bash
curl -u AdminUser:AdminPassword http://localhost:8080/api/v1/get-token
```

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### 📖 List All Books

#### Using Basic Auth:

```bash
curl -u AdminUser:AdminPassword http://localhost:8080/api/v1/books
```

#### Or using JWT:

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/v1/books
```

---

### ➕ Create a Book (JWT Required)

```bash
curl -X POST http://localhost:8080/api/v1/books \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Go Programming",
    "authorList": ["Alan A. A."],
    "publishDate": "2023-01-01",
    "isbn": "123-4567890123"
  }'
```

**Response:**

```json
{
  "uuid": "generated-uuid",
  "name": "Go Programming",
  "authorList": ["Alan A. A."],
  "publishDate": "2023-01-01",
  "isbn": "123-4567890123"
}
```

---

### 🔍 Get Book by UUID (JWT Required)

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/v1/books/<uuid>
```

---

### ✏️ Update a Book (JWT Required)

```bash
curl -X PUT http://localhost:8080/api/v1/books/<uuid> \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Book",
    "authorList": ["New Author"],
    "publishDate": "2025-06-01",
    "isbn": "999-9999999999"
  }'
```

---

### ❌ Delete a Book (JWT Required)

```bash
curl -X DELETE http://localhost:8080/api/v1/books/<uuid> \
  -H "Authorization: Bearer <TOKEN>"
```

---

## 📝 Notes

* All book data is stored in memory and **will be lost on server restart**.
* The JWT secret is hardcoded for demo purposes. Use environment variables for production.
* All requests are logged using Chi's middleware.

```

