package handler

import (
	"net/http"

	"github.com/biswasurmi/book-cli/api/middleware"
	"github.com/biswasurmi/book-cli/service"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router   *chi.Mux
	Handler  *Handler
	Services *service.Services
	Auth     bool
}

func CreateNewServer(h *Handler, services *service.Services, auth bool) *Server {
	r := chi.NewRouter()
	return &Server{
		Router:   r,
		Handler:  h,
		Services: services,
		Auth:     auth,
	}
}

func (s *Server) MountRoutes() {
	s.Router.Post("/api/v1/register", s.Handler.UserHandler.Register)
	s.Router.Post("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		s.Handler.UserHandler.Login(w, r)
	})

	if s.Auth {
		s.Router.Group(func(r chi.Router) {
			r.Use(middleware.BasicAuth(&middleware.BasicAuthConfig{UserService: s.Services.UserService}))
			r.Get("/api/v1/get-token", func(w http.ResponseWriter, r *http.Request) {
				middleware.GetTokenHandler(w, r, s.Auth, s.Services.UserService)
			})
		})
	} else {
		s.Router.Get("/api/v1/get-token", func(w http.ResponseWriter, r *http.Request) {
			middleware.GetTokenHandler(w, r, s.Auth, s.Services.UserService)
		})
	}

	// Protected routes (JWT required when auth=true)
	s.Router.Group(func(r chi.Router) {
		if s.Auth {
			r.Use(middleware.JWTAuth())
		}
		r.Get("/api/v1/books", s.Handler.BookHandler.ListBooks)
		r.Post("/api/v1/books", s.Handler.BookHandler.CreateBook)
		r.Get("/api/v1/books/{uuid}", s.Handler.BookHandler.GetBook)
		r.Put("/api/v1/books/{uuid}", s.Handler.BookHandler.UpdateBook)
		r.Delete("/api/v1/books/{uuid}", s.Handler.BookHandler.DeleteBook)
		r.Get("/api/v1/users/{id}", s.Handler.UserHandler.GetUser)
		r.Get("/api/v1/users/me", s.Handler.UserHandler.GetMe)
		r.Put("/api/v1/users/{id}", s.Handler.UserHandler.UpdateUser)
		r.Delete("/api/v1/users/{id}", s.Handler.UserHandler.Delete)
	})
}