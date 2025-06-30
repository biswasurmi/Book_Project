package test_file

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/biswasurmi/book-cli/api/handler"
	"github.com/biswasurmi/book-cli/domain/entity"
	"github.com/biswasurmi/book-cli/infrastructure/persistance/inmemory"
	"github.com/biswasurmi/book-cli/service"
	"github.com/go-chi/jwtauth/v5"
)

func BasicAuthHeader(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func GenerateJWTToken(tokenAuth *jwtauth.JWTAuth) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": "test-user"})
	return "Bearer " + tokenString
}

func Test_Get_Token(t *testing.T) {
	repo := inmemory.NewBookRepo()
	svc := service.NewBookService(repo)
	bookHandler := handler.NewBookHandler(svc)
	tokenAuth := jwtauth.New("HS256", []byte("supersecretkey123"), nil)
	s := handler.CreateNewServer(bookHandler, true, tokenAuth)
	s.MountRoutes()

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	test := []Test{
		{
			"GET",
			"/api/v1/get-token",
			nil,
			BasicAuthHeader("admin", "admin123"),
			http.StatusOK,
		},
		{
			"GET",
			"/api/v1/get-token",
			nil,
			BasicAuthHeader("wrong", "wrong123"),
			http.StatusUnauthorized,
		},
		{
			"GET",
			"/api/v1/get-token",
			nil,
			"",
			http.StatusUnauthorized,
		},
	}

	for _, i := range test {
		req, _ := http.NewRequest(i.method, i.url, i.body)
		if i.token != "" {
			req.Header.Set("Authorization", i.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, i.expectedStatusCode, response.Code)
	}
}
func Test_AllBookList(t *testing.T) {
	repo := inmemory.NewBookRepo()
	svc := service.NewBookService(repo)
	bookHandler := handler.NewBookHandler(svc)
	tokenAuth := jwtauth.New("HS256", []byte("supersecretkey123"), nil)
	s := handler.CreateNewServer(bookHandler, true, tokenAuth)
	s.MountRoutes()

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	test := []Test{
		{
			"GET",
			"/api/v1/books",
			nil,
			BasicAuthHeader("admin", "admin123"),
			http.StatusOK,
		},
		{
			"GET",
			"/api/v1/books",
			nil,
			BasicAuthHeader("wrong", "wrong123"),
			http.StatusUnauthorized,
		},
		{
			"GET",
			"/api/v1/books",
			nil,
			"",
			http.StatusUnauthorized,
		},
	}

	for _, i := range test {
		req, _ := http.NewRequest(i.method, i.url, i.body)
		if i.token != "" {
			req.Header.Set("Authorization", i.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, i.expectedStatusCode, response.Code)
	}
}

func Test_Create_Book(t *testing.T) {

	repo := inmemory.NewBookRepo()
	svc := service.NewBookService(repo)
	bookHandler := handler.NewBookHandler(svc)
	tokenAuth := jwtauth.New("HS256", []byte("supersecretkey123"), nil)
	s := handler.CreateNewServer(bookHandler, true, tokenAuth)
	s.MountRoutes()

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	test := []Test{
		{
			"POST",
			"/api/v1/books",
			bytes.NewReader([]byte(`{"name":"Learn API","authorList":["Urmi"],"publishDate":"2022-01-02","isbn":"0999-0555-5914"}`)),
			GenerateJWTToken(tokenAuth),
			http.StatusCreated,
		},
		{
			"POST",
			"/api/v1/books",
			bytes.NewReader([]byte(`{"name":"Learn API","authorList":"Urmi","publishDate":"2022-01-02","isbn":"0999-0555-5914"}`)),
			GenerateJWTToken(tokenAuth),
			http.StatusBadRequest,
		},
		{
			"POST",
			"/api/v1/books",
			bytes.NewReader([]byte(`{"name":"Learn API","authorList":["Urmi"],"publishDate":"2022-01-02","isbn":"0999-0555-5914"}`)),
			"Bearer invalid.token.here",
			http.StatusUnauthorized,
		},
		{
			"POST",
			"/api/v1/books",
			bytes.NewReader([]byte(`{"name":"Learn API","authorList":["Urmi"],"publishDate":"2022-01-02","isbn":"0999-0555-5914"}`)),
			"",
			http.StatusUnauthorized,
		},
	}

	for _, i := range test {
		req, _ := http.NewRequest(i.method, i.url, i.body)
		if i.token != "" {
			req.Header.Set("Authorization", i.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, i.expectedStatusCode, response.Code)
	}
}

func Test_Get_Book_With_id(t *testing.T) {

	repo := inmemory.NewBookRepo()
	svc := service.NewBookService(repo)
	bookHandler := handler.NewBookHandler(svc)
	tokenAuth := jwtauth.New("HS256", []byte("supersecretkey123"), nil)
	s := handler.CreateNewServer(bookHandler, true, tokenAuth)
	s.MountRoutes()

	book := entity.Book{
		UUID:        "123e4567-e89b-12d3-a456-426614174001",
		Name:        "Learn API",
		AuthorList:  []string{"Urmi"},
		PublishDate: "2022-01-02",
		ISBN:        "0999-0555-5914",
	}
	repo.CreateBook(book)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	test := []Test{
		{
			"GET",
			"/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			nil,
			GenerateJWTToken(tokenAuth),
			http.StatusOK,
		},
		{
			"GET",
			"/api/v1/books/non-existent-uuid",
			nil,
			GenerateJWTToken(tokenAuth),
			http.StatusNotFound,
		},
		{
			"GET",
			"/api/v1/books/",
			nil,
			GenerateJWTToken(tokenAuth),
			http.StatusNotFound,
		},
		{
			"GET",
			"/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			nil,
			"Bearer invalid.token.here",
			http.StatusUnauthorized,
		},
		{
			"GET",
			"/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			nil,
			"",
			http.StatusUnauthorized,
		},
	}

	for _, i := range test {
		req, _ := http.NewRequest(i.method, i.url, i.body)
		if i.token != "" {
			req.Header.Set("Authorization", i.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, i.expectedStatusCode, response.Code)
	}
}

func Test_Update_Book(t *testing.T) {
	repo := inmemory.NewBookRepo()
	svc := service.NewBookService(repo)
	bookHandler := handler.NewBookHandler(svc)
	tokenAuth := jwtauth.New("HS256", []byte("supersecretkey123"), nil)
	s := handler.CreateNewServer(bookHandler, true, tokenAuth)
	s.MountRoutes()

	book := entity.Book{
		UUID:        "123e4567-e89b-12d3-a456-426614174001",
		Name:        "Learn API",
		AuthorList:  []string{"Urmi"},
		PublishDate: "2022-01-02",
		ISBN:        "0999-0555-5914",
	}
	repo.CreateBook(book)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	test := []Test{
		{
			"PUT",
			"/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			bytes.NewReader([]byte(`{"name":"Updated API","authorList":["Biswas"],"publishDate":"2023-01-02","isbn":"0999-0555-5954"}`)),
			GenerateJWTToken(tokenAuth),
			http.StatusOK,
		},
		{
			"PUT",
			"/api/v1/books/non-existent-uuid",
			bytes.NewReader([]byte(`{"name":"Updated API","authorList":["Biswas"],"publishDate":"2023-01-02","isbn":"0999-0555-5954"}`)),
			GenerateJWTToken(tokenAuth),
			http.StatusNotFound,
		},
		{
			"PUT",
			"/api/v1/books/",
			bytes.NewReader([]byte(`{"name":"Updated API","authorList":["Biswas"],"publishDate":"2023-01-02","isbn":"0999-0555-5954"}`)),
			GenerateJWTToken(tokenAuth),
			http.StatusNotFound,
		},
		{
			"PUT",
			"/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			bytes.NewReader([]byte(`{"name":"Updated API","authorList":"Biswas","publishDate":"2023-01-02","isbn":"0999-0555-5954"}`)),
			GenerateJWTToken(tokenAuth),
			http.StatusBadRequest,
		},
		{
			"PUT",
			"/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			bytes.NewReader([]byte(`{"name":"Updated API","authorList":["Biswas"],"publishDate":"2023-01-02","isbn":"0999-0555-5954"}`)),
			"Bearer invalid.token.here",
			http.StatusUnauthorized,
		},
		{
			"PUT",
			"/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			bytes.NewReader([]byte(`{"name":"Updated API","authorList":["Biswas"],"publishDate":"2023-01-02","isbn":"0999-0555-5954"}`)),
			"",
			http.StatusUnauthorized,
		},
	}

	for _, i := range test {
		req, _ := http.NewRequest(i.method, i.url, i.body)
		if i.token != "" {
			req.Header.Set("Authorization", i.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, i.expectedStatusCode, response.Code)
	}
}

func Test_Delete_Book(t *testing.T) {
	repo := inmemory.NewBookRepo()
	svc := service.NewBookService(repo)
	bookHandler := handler.NewBookHandler(svc)
	tokenAuth := jwtauth.New("HS256", []byte("supersecretkey123"), nil)
	s := handler.CreateNewServer(bookHandler, true, tokenAuth)
	s.MountRoutes()

	book := entity.Book{
		UUID:        "123e4567-e89b-12d3-a456-426614174001",
		Name:        "Learn API",
		AuthorList:  []string{"Urmi"},
		PublishDate: "2022-01-02",
		ISBN:        "0999-0555-5914",
	}
	repo.CreateBook(book)

	type Test struct {
		method             string
		url                string
		body               io.Reader
		token              string
		expectedStatusCode int
	}

	test := []Test{
		{
			"DELETE",
			"/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			nil,
			GenerateJWTToken(tokenAuth),
			http.StatusNoContent,
		},
		{
			"DELETE",
			"/api/v1/books/non-existent-uuid",
			nil,
			GenerateJWTToken(tokenAuth),
			http.StatusNotFound,
		},
		{
			"DELETE",
			"/api/v1/books/",
			nil,
			GenerateJWTToken(tokenAuth),
			http.StatusNotFound,
		},
		{
			"DELETE",
			"/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			nil,
			"Bearer invalid.token.here",
			http.StatusUnauthorized,
		},
		{
			"DELETE",
			"/api/v1/books/123e4567-e89b-12d3-a456-426614174001",
			nil,
			"",
			http.StatusUnauthorized,
		},
	}

	for _, i := range test {
		req, _ := http.NewRequest(i.method, i.url, i.body)
		if i.token != "" {
			req.Header.Set("Authorization", i.token)
		}
		response := executeRequest(req, s)
		checkResponseCode(t, i.expectedStatusCode, response.Code)
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
