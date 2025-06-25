package repository

import (
	"encoding/json"
	"errors"
	"gou-pc/internal/api/model"
	"os"
)

type UserRepository interface {
	GetAll() ([]model.User, error)
	SaveAll([]model.User) error
	FindByID(id string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
}

type fileUserRepository struct {
	file string
}

func NewFileUserRepository(file string) UserRepository {
	return &fileUserRepository{file: file}
}

func (r *fileUserRepository) GetAll() ([]model.User, error) {
	f, err := os.ReadFile(r.file)
	if err != nil {
		return nil, err
	}
	var users []model.User
	if err := json.Unmarshal(f, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *fileUserRepository) SaveAll(users []model.User) error {
	b, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.file, b, 0644)
}

func (r *fileUserRepository) FindByID(id string) (*model.User, error) {
	users, err := r.GetAll()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *fileUserRepository) FindByUsername(username string) (*model.User, error) {
	users, err := r.GetAll()
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}
