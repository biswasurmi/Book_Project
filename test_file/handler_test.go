package test_file

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/biswasurmi/book-cli/api/handler"
	"github.com/biswasurmi/book-cli/domain/entity"
	"github.com/biswasurmi/book-cli/domain/repository"
	"github.com/biswasurmi/book-cli/infrastructure/persistance/inmemory"
	"github.com/biswasurmi/book-cli/service"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

// Load environment variables once during initialization
func init() {
	dir, _ := os.Getwd()
	fmt.Println("Current working directory:", dir)
	if err := godotenv.Load("/mnt/c/Users/Asus/Desktop/Book_Project/.env"); err != nil {
		fmt.Println("Warning: Could not load .env file:", err)
		os.Setenv("JWT_SECRET", "bolaJabeNah")
	}
}

func BasicAuthHeader(email, password string) string {
	auth := email + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func GenerateJWTToken(userID int64) string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "bolaJabeNah" // Fallback for testing
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "" // Avoid panic
	}
	return "Bearer " + tokenString
}

func setupServer(t *testing.T) (*handler.Server, *repository.Repositories) {
	repos := inmemory.GetRepositories()
	services := service.GetServices(repos)
	handlers := &handler.Handler{
		UserHandler: handler.NewUserHandler(services.UserService),
		BookHandler: handler.NewBookHandler(services.BookService),
	}
	s := handler.CreateNewServer(handlers, services, true)
	s.MountRoutes()
	return s, repos
}

func Test_Register(t *testing.T) {
	s, _ := setupServer(t)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	tests := []Test{
		{
			method:             "POST",
			url:                "/api/v1/register",
			body:               bytes.NewReader([]byte(`{"email":"test@example.com","password":"password123"}`)),
			token:              "",
			expectedStatusCode: http.StatusCreated,
		},
		{
			method:             "POST",
			url:                "/api/v1/register",
			body:               bytes.NewReader([]byte(`{"email":"invalid","password":""}`)),
			token:              "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			method:             "POST",
			url:                "/api/v1/register",
			body:               bytes.NewReader([]byte(`{"email":"test2@example.com","password":"short"}`)),
			token:              "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			method:             "POST",
			url:                "/api/v1/register",
			body:               bytes.NewReader([]byte(`not a json`)),
			token:              "",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.url, test.body)
		if test.token != "" {
			req.Header.Set("Authorization", test.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, test.expectedStatusCode, response.Code)
	}
}

func Test_Login(t *testing.T) {
	s, repos := setupServer(t)

	// Pre-create a user
	user := entity.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "$2a$10$bxCN.KcstTAU5I1zkZNe/OYrwD5gUc93lNl5pTit40/ZugB9YwuT6", // Hashed "password123"
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repos.UserRepository.CreateUser(user)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	tests := []Test{
		{
			method:             "POST",
			url:                "/api/v1/login",
			body:               bytes.NewReader([]byte(`{"email":"test@example.com","password":"password123"}`)),
			token:              "",
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             "POST",
			url:                "/api/v1/login",
			body:               bytes.NewReader([]byte(`{"email":"test@example.com","password":"wrong"}`)),
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			method:             "POST",
			url:                "/api/v1/login",
			body:               bytes.NewReader([]byte(`{"email":"nonexistent@example.com","password":"password123"}`)),
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			method:             "POST",
			url:                "/api/v1/login",
			body:               bytes.NewReader([]byte(`not a json`)),
			token:              "",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.url, test.body)
		if test.token != "" {
			req.Header.Set("Authorization", test.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, test.expectedStatusCode, response.Code)
	}
}

func Test_Get_User(t *testing.T) {
	s, repos := setupServer(t)

	// Pre-create a user
	user := entity.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "$2a$10$bxCN.KcstTAU5I1zkZNe/OYrwD5gUc93lNl5pTit40/ZugB9YwuT6",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repos.UserRepository.CreateUser(user)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	tests := []Test{
		{
			method:             "GET",
			url:                "/api/v1/users/1",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             "GET",
			url:                "/api/v1/users/999",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			method:             "GET",
			url:                "/api/v1/users/invalid",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			method:             "GET",
			url:                "/api/v1/users/1",
			body:               nil,
			token:              "Bearer invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			method:             "GET",
			url:                "/api/v1/users/1",
			body:               nil,
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.url, test.body)
		if test.token != "" {
			req.Header.Set("Authorization", test.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, test.expectedStatusCode, response.Code)
	}
}

func Test_Get_Me(t *testing.T) {
	s, repos := setupServer(t)

	// Pre-create a user
	user := entity.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "$2a$10$bxCN.KcstTAU5I1zkZNe/OYrwD5gUc93lNl5pTit40/ZugB9YwuT6",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repos.UserRepository.CreateUser(user)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	tests := []Test{
		{
			method:             "GET",
			url:                "/api/v1/users/me",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             "GET",
			url:                "/api/v1/users/me",
			body:               nil,
			token:              GenerateJWTToken(999), // Non-existent user
			expectedStatusCode: http.StatusNotFound,
		},
		{
			method:             "GET",
			url:                "/api/v1/users/me",
			body:               nil,
			token:              "Bearer invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			method:             "GET",
			url:                "/api/v1/users/me",
			body:               nil,
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.url, test.body)
		if test.token != "" {
			req.Header.Set("Authorization", test.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, test.expectedStatusCode, response.Code)
	}
}

func Test_Update_User(t *testing.T) {
	s, repos := setupServer(t)

	// Pre-create a user
	user := entity.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "$2a$10$bxCN.KcstTAU5I1zkZNe/OYrwD5gUc93lNl5pTit40/ZugB9YwuT6",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repos.UserRepository.CreateUser(user)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	tests := []Test{
		{
			method:             "PUT",
			url:                "/api/v1/users/1",
			body:               bytes.NewReader([]byte(`{"email":"updated@example.com","password":"newpassword123","username":"testuser"}`)),
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             "PUT",
			url:                "/api/v1/users/999",
			body:               bytes.NewReader([]byte(`{"email":"updated@example.com","password":"newpassword123","username":"testuser"}`)),
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			method:             "PUT",
			url:                "/api/v1/users/invalid",
			body:               bytes.NewReader([]byte(`{"email":"updated@example.com","password":"newpassword123","username":"testuser"}`)),
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			method:             "PUT",
			url:                "/api/v1/users/1",
			body:               bytes.NewReader([]byte(`{"email":"","password":"newpassword123"}`)),
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			method:             "PUT",
			url:                "/api/v1/users/1",
			body:               bytes.NewReader([]byte(`not a json`)),
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			method:             "PUT",
			url:                "/api/v1/users/1",
			body:               bytes.NewReader([]byte(`{"email":"updated@example.com","password":"newpassword123","username":"testuser"}`)),
			token:              "Bearer invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			method:             "PUT",
			url:                "/api/v1/users/1",
			body:               bytes.NewReader([]byte(`{"email":"updated@example.com","password":"newpassword123","username":"testuser"}`)),
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.url, test.body)
		if test.token != "" {
			req.Header.Set("Authorization", test.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, test.expectedStatusCode, response.Code)
	}
}

func Test_Delete_User(t *testing.T) {
	s, repos := setupServer(t)

	// Pre-create a user
	user := entity.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "$2a$10$bxCN.KcstTAU5I1zkZNe/OYrwD5gUc93lNl5pTit40/ZugB9YwuT6",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repos.UserRepository.CreateUser(user)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	tests := []Test{
		{
			method:             "DELETE",
			url:                "/api/v1/users/1",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusNoContent,
		},
		{
			method:             "DELETE",
			url:                "/api/v1/users/999",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			method:             "DELETE",
			url:                "/api/v1/users/invalid",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			method:             "DELETE",
			url:                "/api/v1/users/1",
			body:               nil,
			token:              "Bearer invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			method:             "DELETE",
			url:                "/api/v1/users/1",
			body:               nil,
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.url, test.body)
		if test.token != "" {
			req.Header.Set("Authorization", test.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, test.expectedStatusCode, response.Code)
	}
}

func Test_AllBookList(t *testing.T) {
	s, repos := setupServer(t)

	// Pre-create a user for JWT
	user := entity.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "$2a$10$bxCN.KcstTAU5I1zkZNe/OYrwD5gUc93lNl5pTit40/ZugB9YwuT6",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repos.UserRepository.CreateUser(user)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	tests := []Test{
		{
			method:             "GET",
			url:                "/api/v1/books",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             "GET",
			url:                "/api/v1/books",
			body:               nil,
			token:              "Bearer invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			method:             "GET",
			url:                "/api/v1/books",
			body:               nil,
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.url, test.body)
		if test.token != "" {
			req.Header.Set("Authorization", test.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, test.expectedStatusCode, response.Code)
	}
}

func Test_Create_Book(t *testing.T) {
	s, repos := setupServer(t)

	// Pre-create a user for JWT
	user := entity.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "$2a$10$bxCN.KcstTAU5I1zkZNe/OYrwD5gUc93lNl5pTit40/ZugB9YwuT6",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repos.UserRepository.CreateUser(user)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	tests := []Test{
		{
			method:             "POST",
			url:                "/api/v1/books",
			body:               bytes.NewReader([]byte(`{"name":"Learn API","authorList":["Urmi"],"publishDate":"2022-01-02","isbn":"0999-0555-5914"}`)),
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusCreated,
		},
		{
			method:             "POST",
			url:                "/api/v1/books",
			body:               bytes.NewReader([]byte(`{"name":"Learn API","authorList":"Urmi","publishDate":"2022-01-02","isbn":"0999-0555-5914"}`)),
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			method:             "POST",
			url:                "/api/v1/books",
			body:               bytes.NewReader([]byte(`{"name":"Learn API","authorList":["Urmi"],"publishDate":"2022-01-02","isbn":"0999-0555-5914"}`)),
			token:              "Bearer invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			method:             "POST",
			url:                "/api/v1/books",
			body:               bytes.NewReader([]byte(`{"name":"Learn API","authorList":["Urmi"],"publishDate":"2022-01-02","isbn":"0999-0555-5914"}`)),
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.url, test.body)
		if test.token != "" {
			req.Header.Set("Authorization", test.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, test.expectedStatusCode, response.Code)
	}
}

func Test_Get_Book_With_id(t *testing.T) {
	s, repos := setupServer(t)

	// Pre-create a user for JWT
	user := entity.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "$2a$10$bxCN.KcstTAU5I1zkZNe/OYrwD5gUc93lNl5pTit40/ZugB9YwuT6",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repos.UserRepository.CreateUser(user)

	// Pre-create a book
	book := entity.Book{
		UUID:        "123e4567-e89b-12d3-a456-426614174001",
		Name:        "Learn API",
		AuthorList:  []string{"Urmi"},
		PublishDate: "2022-01-02",
		ISBN:        "0999-0555-5914",
	}
	repos.BookRepository.CreateBook(book)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	tests := []Test{
		{
			method:             "GET",
			url:                "/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             "GET",
			url:                "/api/v1/books/non-existent-uuid",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			method:             "GET",
			url:                "/api/v1/books/",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			method:             "GET",
			url:                "/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			body:               nil,
			token:              "Bearer invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			method:             "GET",
			url:                "/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			body:               nil,
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.url, test.body)
		if test.token != "" {
			req.Header.Set("Authorization", test.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, test.expectedStatusCode, response.Code)
	}
}

func Test_Update_Book(t *testing.T) {
	s, repos := setupServer(t)

	// Pre-create a user for JWT
	user := entity.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "$2a$10$bxCN.KcstTAU5I1zkZNe/OYrwD5gUc93lNl5pTit40/ZugB9YwuT6",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repos.UserRepository.CreateUser(user)

	// Pre-create a book
	book := entity.Book{
		UUID:        "123e4567-e89b-12d3-a456-426614174001",
		Name:        "Learn API",
		AuthorList:  []string{"Urmi"},
		PublishDate: "2022-01-02",
		ISBN:        "0999-0555-5914",
	}
	repos.BookRepository.CreateBook(book)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	tests := []Test{
		{
			method:             "PUT",
			url:                "/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			body:               bytes.NewReader([]byte(`{"name":"Updated API","authorList":["Biswas"],"publishDate":"2023-01-02","isbn":"0999-0555-5954"}`)),
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             "PUT",
			url:                "/api/v1/books/non-existent-uuid",
			body:               bytes.NewReader([]byte(`{"name":"Updated API","authorList":["Biswas"],"publishDate":"2023-01-02","isbn":"0999-0555-5954"}`)),
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			method:             "PUT",
			url:                "/api/v1/books/",
			body:               bytes.NewReader([]byte(`{"name":"Updated API","authorList":["Biswas"],"publishDate":"2023-01-02","isbn":"0999-0555-5954"}`)),
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			method:             "PUT",
			url:                "/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			body:               bytes.NewReader([]byte(`{"name":"Updated API","authorList":"Biswas","publishDate":"2023-01-02","isbn":"0999-0555-5954"}`)),
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			method:             "PUT",
			url:                "/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			body:               bytes.NewReader([]byte(`{"name":"Updated API","authorList":["Biswas"],"publishDate":"2023-01-02","isbn":"0999-0555-5954"}`)),
			token:              "Bearer invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			method:             "PUT",
			url:                "/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			body:               bytes.NewReader([]byte(`{"name":"Updated API","authorList":["Biswas"],"publishDate":"2023-01-02","isbn":"0999-0555-5954"}`)),
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.url, test.body)
		if test.token != "" {
			req.Header.Set("Authorization", test.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, test.expectedStatusCode, response.Code)
	}
}

func Test_Delete_Book(t *testing.T) {
	s, repos := setupServer(t)

	// Pre-create a user for JWT
	user := entity.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "$2a$10$bxCN.KcstTAU5I1zkZNe/OYrwD5gUc93lNl5pTit40/ZugB9YwuT6",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repos.UserRepository.CreateUser(user)

	// Pre-create a book
	book := entity.Book{
		UUID:        "123e4567-e89b-12d3-a456-426614174001",
		Name:        "Learn API",
		AuthorList:  []string{"Urmi"},
		PublishDate: "2022-01-02",
		ISBN:        "0999-0555-5914",
	}
	repos.BookRepository.CreateBook(book)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	tests := []Test{
		{
			method:             "DELETE",
			url:                "/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusNoContent,
		},
		{
			method:             "DELETE",
			url:                "/api/v1/books/non-existent-uuid",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			method:             "DELETE",
			url:                "/api/v1/books/",
			body:               nil,
			token:              GenerateJWTToken(1),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			method:             "DELETE",
			url:                "/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			body:               nil,
			token:              "Bearer invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			method:             "DELETE",
			url:                "/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			body:               nil,
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.url, test.body)
		if test.token != "" {
			req.Header.Set("Authorization", test.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, test.expectedStatusCode, response.Code)
	}
}

func executeRequest(req *http.Request, s *handler.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}


// go test -v ./test_file