package handler

import "github.com/biswasurmi/book-cli/service"

type Handler struct {
	BookHandler *BookHandler
	UserHandler *UserHandler
}

func GetHandlers(services *service.Services) *Handler {
	return &Handler{
		BookHandler: NewBookHandler(services.BookService),
		UserHandler: NewUserHandler(services.UserService, nil), // tokenAuth will be passed in cmd
	}
}