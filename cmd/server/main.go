package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/di"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/handler"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/middleware"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/pkg/database"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize dependency container
	container, err := di.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	// Set up Gin
	if !cfg.Logging.IsDevelopment {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middlewares
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger(container.Logger))

	// Register health check routes
	healthHandler := handler.NewHealthHandler(container.DB, container.Logger)
	healthHandler.RegisterRoutes(router)

	// API routes
	api := router.Group("/api/v1")

	// User routes (protected by authentication)
	userRoutes := api.Group("/user")
	userRoutes.Use(middleware.Authentication(container.Logger))

	// Register profile handler routes
	userProfileHandler := handler.NewUserProfileHandler(container.UserProfileService, container.Logger)
	userProfileHandler.RegisterRoutes(userRoutes)

	// Run database migrations
	if cfg.Database.RunMigrations {
		container.Logger.Info("Running database migrations")
		if err := database.RunMigrations(cfg.Database.DSN()); err != nil {
			container.Logger.Fatal("Failed to run database migrations", zap.Error(err))
		}
		container.Logger.Info("Database migrations completed")
	}

	// Start the server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.RequestTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.RequestTimeoutSec) * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Run the server in a goroutine
	go func() {
		container.Logger.Info("Starting server", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			container.Logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Logger.Info("Shutting down server...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		container.Logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	container.Logger.Info("Server exited properly")
}
