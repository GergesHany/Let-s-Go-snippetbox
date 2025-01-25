package mocks

import (
	"snippetbox.alexedwards.net/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, hashedPassword string) error {
	if email == "dupe@example.com" {
		return models.ErrDuplicateEmail
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
    	return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	if id == 1 {
		return true, nil
	}
	return false, nil
}