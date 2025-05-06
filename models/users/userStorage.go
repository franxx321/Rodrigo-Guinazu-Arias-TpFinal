package users

import "errors"

var ErrNotFound = errors.New("user not found")
var ErrEmptyID = errors.New("empty user ID")

type UserStorage interface {
	Set(user *User) error
	Read(id string) (*User, error)
	Delete(id string) error
}
type LocalStorage struct {
	m map[string]*User
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		m: map[string]*User{},
	}
}

func (l *LocalStorage) Set(user *User) error {
	if user.ID == "" {
		return ErrEmptyID
	}

	l.m[user.ID] = user
	return nil
}

func (l *LocalStorage) Read(id string) (*User, error) {
	u, ok := l.m[id]
	if !ok {
		return nil, ErrNotFound
	}

	return u, nil
}

func (l *LocalStorage) Delete(id string) error {
	_, err := l.Read(id)
	if err != nil {
		return err
	}

	delete(l.m, id)
	return nil
}
