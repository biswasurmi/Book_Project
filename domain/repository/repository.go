package repository



// Repositories aggregates all repository interfaces
type Repositories struct {
	BookRepository BookRepository
	UserRepository UserRepository
}