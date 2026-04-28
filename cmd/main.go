package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/codelogydev/core-go/cache"
	"github.com/codelogydev/core-go/logger"
	coreMiddleware "github.com/codelogydev/core-go/middleware"
	"github.com/codelogydev/core-go/storage"

	"github.com/codelogydev/template-go-api/internal/database"
	"github.com/codelogydev/template-go-api/internal/handler"
	"github.com/codelogydev/template-go-api/internal/repository"
	"github.com/codelogydev/template-go-api/internal/service"
)

func main() {
	_ = godotenv.Load()

	logger.Init()
	defer logger.Log.Sync()

	if err := database.Connect(); err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}

	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		if err := cache.Init(redisURL); err != nil {
			logger.Log.Warn("redis unavailable, cache disabled", zap.Error(err))
		}
	}

	if endpoint := os.Getenv("STORAGE_ENDPOINT"); endpoint != "" {
		if err := storage.Init(
			endpoint,
			os.Getenv("STORAGE_ACCESS_KEY"),
			os.Getenv("STORAGE_SECRET_KEY"),
			os.Getenv("STORAGE_USE_SSL") == "true",
		); err != nil {
			logger.Log.Warn("storage unavailable", zap.Error(err))
		}
	}

	userRepo := repository.NewUserRepository(database.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := gin.New()
	r.Use(coreMiddleware.Recovery())
	r.Use(coreMiddleware.Logger())

	r.GET("/", handler.HealthCheck)
	r.GET("/users", userHandler.GetUsers)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("server exited cleanly")
}

