package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	coreMiddleware "github.com/codelogydev/core-go/middleware"
	"github.com/codelogydev/core-go/logger"
	"github.com/codelogydev/template-go-api/internal/database"
	"github.com/codelogydev/template-go-api/internal/handler"
	"github.com/codelogydev/template-go-api/internal/repository"
	"github.com/codelogydev/template-go-api/internal/service"
)

func main() {
	_ = godotenv.Load()

	logger.Init()
	defer logger.Log.Sync()

	err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := gin.New()
	r.Use(coreMiddleware.Recovery())
	r.Use(coreMiddleware.Logger())

	r.GET("/", handler.HealthCheck)
	r.GET("/users", userHandler.GetUsers)

	r.Run(":8080")
}
