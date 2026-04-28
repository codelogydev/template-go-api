package service

import (
	"context"
	"errors"

	"github.com/codelogydev/core-go/auth"
	"github.com/codelogydev/template-go-api/internal/dto"
	"github.com/codelogydev/template-go-api/internal/model"
	"github.com/codelogydev/template-go-api/internal/repository"
)

type UserService interface {
	GetAllUsers() ([]dto.UserResponse, error)
	LoginWithGoogle(googleID, email, name string) (*dto.LoginResponse, error)
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

func (s *userService) LoginWithGoogle(googleID, email, name string) (*dto.LoginResponse, error) {
	ctx := context.Background()

	user, err := s.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrNotFound) {
		user, err = s.repo.Create(ctx, name, email)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  toUserResponse(*user),
	}, nil
}

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

