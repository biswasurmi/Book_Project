package handler

import (
	"net/http"
	"github.com/biswasurmi/book-cli/api/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

type Server struct {
	Router      *chi.Mux
	BookHandler *BookHandler
	AuthEnabled bool
	TokenAuth   *jwtauth.JWTAuth
}

func CreateNewServer(bookHandler *BookHandler, authEnabled bool, tokenAuth *jwtauth.JWTAuth) *Server {
	return &Server{
		Router:      chi.NewRouter(),
		BookHandler: bookHandler,
		AuthEnabled: authEnabled,
		TokenAuth:   tokenAuth,
	}
}

func (s *Server) MountRoutes() {
	s.Router.Use(chiMiddleware.Logger)

	s.Router.Group(func(r chi.Router) {
		if s.AuthEnabled {
			r.Use(middleware.BasicAuth)
		}
		r.Get("/api/v1/books", s.BookHandler.ListBooks)
	})

	s.Router.Group(func(r chi.Router) {
		if s.AuthEnabled {
			r.Use(middleware.JWTAuth(s.TokenAuth))
		}
		r.Post("/api/v1/books", s.BookHandler.CreateBook)
		r.Get("/api/v1/books/{uuid}", s.BookHandler.GetBook)
		r.Put("/api/v1/books/{uuid}", s.BookHandler.UpdateBook)
		r.Delete("/api/v1/books/{uuid}", s.BookHandler.DeleteBook)
	})

	s.Router.Get("/api/v1/get-token", func(w http.ResponseWriter, r *http.Request) {
		middleware.GetTokenHandler(w, r, s.TokenAuth, s.AuthEnabled)
	})
}