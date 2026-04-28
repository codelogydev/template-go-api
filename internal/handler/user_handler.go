package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/codelogydev/core-go/response"
	"github.com/codelogydev/template-go-api/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, users)
}
