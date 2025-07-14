package service

import (
	"github.com/biswasurmi/book-cli/domain/entity"
	"github.com/biswasurmi/book-cli/domain/repository"
)

type BookService interface {
	ListBooks() ([]entity.Book, error)
	CreateBook(book entity.Book) (entity.Book, error)
	GetBook(uuid string) (entity.Book, error)
	UpdateBook(book entity.Book) (entity.Book, error)
	DeleteBook(uuid string) error
}

type bookService struct {
	bookRepo repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{bookRepo: bookRepo}
}

func (s *bookService) ListBooks() ([]entity.Book, error) {
	return s.bookRepo.GetAllBooks()
}

func (s *bookService) CreateBook(book entity.Book) (entity.Book, error) {
	return s.bookRepo.CreateBook(book)
}

func (s *bookService) GetBook(uuid string) (entity.Book, error) {
	return s.bookRepo.GetBook(uuid)
}

func (s *bookService) UpdateBook(book entity.Book) (entity.Book, error) {
	return s.bookRepo.UpdateBook(book)
}

func (s *bookService) DeleteBook(uuid string) error {
	return s.bookRepo.DeleteBook(uuid)
}