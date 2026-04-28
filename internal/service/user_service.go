package service

import (
	"github.com/codelogydev/template-go-api/internal/dto"
	"github.com/codelogydev/template-go-api/internal/model"
	"github.com/codelogydev/template-go-api/internal/repository"
)

type UserService interface {
	GetAllUsers() ([]dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return toUserResponseList(users), nil
}

// toUserResponse memetakan model.User ke dto.UserResponse.
// Ini memastikan field sensitif (PasswordHash, dll) tidak pernah keluar ke client.
func toUserResponse(u model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func toUserResponseList(users []model.User) []dto.UserResponse {
	result := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		result = append(result, toUserResponse(u))
	}
	return result
}
