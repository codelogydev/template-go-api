package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/codelogydev/core-go/response"
)

func HealthCheck(c *gin.Context) {
	response.Success(c, "API is running 🚀")
}
