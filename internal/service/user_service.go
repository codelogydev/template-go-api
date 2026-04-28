package service

import (
	"github.com/codelogydev/template-go-api/internal/model"
	"github.com/codelogydev/template-go-api/internal/repository"
)

type UserService interface {
	GetAllUsers() ([]model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.repo.GetAll()
}
