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
	s.userStorage.Set(user)
	return nil
}

func (s *UserService) Get(id string) (*User, error) {
	return s.userStorage.Read(id)
}

func (s *UserService) Update(id string, user *UpdateFields) (*User, error) {
	existing, err := s.userStorage.Read(id)
	if err != nil {
		return nil, err
	}

	if user.Name != nil {
		existing.Name = *user.Name
	}

	if user.Address != nil {
		existing.Address = *user.Address
	}

	if user.NickName != nil {
		existing.NickName = *user.NickName
	}

	existing.UpdatedAt = time.Now()
	existing.Version++

	if err := s.userStorage.Set(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *UserService) Delete(id string) error {
	return s.userStorage.Delete(id)
}
