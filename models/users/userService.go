package users

import (
	"github.com/google/uuid"
	"time"
)

type UserService struct {
	userStorage UserStorage
}

func NewUserService(userStorage UserStorage) *UserService {
	return &UserService{
		userStorage: userStorage,
	}
}

func (s *UserService) Create(user *User) error {
	user.ID = uuid.NewString()
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.Version = 1

	return nil
}

func (s *UserService) Get(id string) (*User, error) {
	return s.userStorage.Read(id)
}

func (s *UserService) Update(id string, name, address, nickname string) (*User, error) {
	user, err := s.userStorage.Read(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		user.Name = name
	}
	if address != "" {
		user.Address = address
	}
	if nickname != "" {
		user.NickName = nickname
	}

	user.UpdatedAt = time.Now()
	user.Version++

	err = s.userStorage.Set(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(id string) error {
	return s.userStorage.Delete(id)
}
