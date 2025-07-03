package service

import "github.com/biswasurmi/book-cli/domain/repository"

type Services struct {
	BookService BookService
	UserService UserService
}

func GetServices(repos *repository.Repositories) *Services {
	return &Services{
		BookService: NewBookService(repos.BookRepository),
		UserService: NewUserService(repos.UserRepository),
	}
}