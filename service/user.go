package service

import (
    "github.com/biswasurmi/book-cli/domain/entity"
    "github.com/biswasurmi/book-cli/domain/repository"
)

type UserService interface {
    CreateUser(user entity.User) (entity.User, error)
    GetByID(id int64) (entity.User, error)
    GetByEmail(email string) (entity.User, error)
    Update(user entity.User) (entity.User, error)
    Delete(user entity.User) error
    Authenticate(email, password string) (entity.User, error)
}

type userService struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
    return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(user entity.User) (entity.User, error) {
    return s.userRepo.CreateUser(user)
}

func (s *userService) GetByID(id int64) (entity.User, error) {
    return s.userRepo.GetByID(id)
}

func (s *userService) GetByEmail(email string) (entity.User, error) {
    return s.userRepo.GetByEmail(email)
}

func (s *userService) Update(user entity.User) (entity.User, error) {
    return s.userRepo.Update(user)
}

func (s *userService) Delete(user entity.User) error {
    return s.userRepo.Delete(user.ID)
}

func (s *userService) Authenticate(email, password string) (entity.User, error) {
    return s.userRepo.Authenticate(email, password)
}