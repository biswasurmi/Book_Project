package inmemory

import (
	"errors"
	"time"

	"github.com/biswasurmi/book-cli/domain/entity"
	"github.com/biswasurmi/book-cli/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type userRepo struct {
	users map[int64]entity.User
}

func NewUserRepo() repository.UserRepository {
	return &userRepo{
		users: make(map[int64]entity.User),
	}
}

func (r *userRepo) CreateUser(user entity.User) (entity.User, error) {
	r.users[user.ID] = user
	return user, nil
}

func (r *userRepo) GetByID(id int64) (entity.User, error) {
	user, exists := r.users[id]
	if !exists {
		return entity.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepo) GetByEmail(email string) (entity.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return entity.User{}, errors.New("user not found")
}

func (r *userRepo) Update(user entity.User) (entity.User, error) {
	if _, exists := r.users[user.ID]; !exists {
		return entity.User{}, errors.New("user not found")
	}
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return user, nil
}

func (r *userRepo) Delete(id int64) error {
	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}

func (r *userRepo) Authenticate(email, password string) (entity.User, error) {
	
	user, err := r.GetByEmail(email)
	if err != nil {
		return entity.User{}, errors.New("invalid credentials")
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return entity.User{}, errors.New("invalid credentials")
	}
	return user, nil
}