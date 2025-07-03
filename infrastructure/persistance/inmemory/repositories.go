package inmemory

import (
    "github.com/biswasurmi/book-cli/domain/repository"
)

// GetRepositories returns a *repository.Repositories with in-memory implementations.
func GetRepositories() *repository.Repositories {
    return &repository.Repositories{
        BookRepository: NewBookRepo(),
        UserRepository: NewUserRepo(),
    }
}