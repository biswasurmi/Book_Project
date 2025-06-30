package inmemory

import (
	"errors"
	"github.com/biswasurmi/book-cli/domain/entity"
	"github.com/biswasurmi/book-cli/domain/repository"
)

type bookRepo struct {
	books map[string]entity.Book
}

func NewBookRepo() repository.BookRepository {
	return &bookRepo{
		books: make(map[string]entity.Book),
	}
}

func (b *bookRepo) GetAllBooks() ([]entity.Book, error) {
	var result []entity.Book
	for _, book := range b.books {
		result = append(result, book)
	}
	return result, nil
}

func (b *bookRepo) CreateBook(book entity.Book) (entity.Book, error) {
	b.books[book.UUID] = book
	return book, nil
}

func (b *bookRepo) GetBook(uuid string) (entity.Book, error) {
	book, exists := b.books[uuid]
	if !exists {
		return entity.Book{}, errors.New("book not found")
	}
	return book, nil
}

func (b *bookRepo) UpdateBook(book entity.Book) (entity.Book, error) {
	if _, exists := b.books[book.UUID]; !exists {
		return entity.Book{}, errors.New("book not found")
	}
	b.books[book.UUID] = book
	return book, nil
}

func (b *bookRepo) DeleteBook(uuid string) error {
	if _, exists := b.books[uuid]; !exists {
		return errors.New("book not found")
	}
	delete(b.books, uuid)
	return nil
}