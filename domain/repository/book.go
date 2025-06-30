package repository

import "github.com/biswasurmi/book-cli/domain/entity"

type BookRepository interface {
	GetAllBooks() ([]entity.Book, error)
	CreateBook(book entity.Book) (entity.Book, error)
	GetBook(uuid string) (entity.Book, error)
	UpdateBook(book entity.Book) (entity.Book, error)
	DeleteBook(uuid string) error
}