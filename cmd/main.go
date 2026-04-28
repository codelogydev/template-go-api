package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/codelogydev/template-go-api/internal/database"
	"github.com/codelogydev/template-go-api/internal/handler"
	"github.com/codelogydev/template-go-api/internal/repository"
	"github.com/codelogydev/template-go-api/internal/service"
)

func main() {
	_ = godotenv.Load()

	err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Dependency injection: DB → repo → service → handler
	userRepo := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()

	r.GET("/", handler.HealthCheck)
	r.GET("/users", userHandler.GetUsers)

	r.Run(":8080")
}
