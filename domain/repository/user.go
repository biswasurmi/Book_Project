package repository

import "github.com/biswasurmi/book-cli/domain/entity"



type UserRepository interface {
	CreateUser(user entity.User) (entity.User, error)
	GetByID (id int64) (entity.User, error)
	GetByEmail (email string) (entity.User, error)
	Update(user entity.User) (entity.User, error)
	Delete(id int64) error
	Authenticate(email, password string) (entity.User, error)
}