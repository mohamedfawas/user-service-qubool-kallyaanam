// user-service-qubool-kallyaanam/cmd/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// MongoDB connection
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://mongo:27017/user_db" // Default in Docker
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v", err)
	} else {
		log.Println("Successfully connected to MongoDB")
		defer client.Disconnect(ctx)
	}

	router := gin.Default()

	// Health Check endpoint
	router.GET("/health", func(c *gin.Context) {
		// Check MongoDB health
		mongoStatus := "UP"
		if client != nil {
			err := client.Ping(ctx, readpref.Primary())
			if err != nil {
				mongoStatus = "DOWN"
			}
		} else {
			mongoStatus = "DOWN"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"service": "user-service",
			"version": "0.1.0",
			"mongodb": mongoStatus,
		})
	})

	// Start server
	srv := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a timeout of 5 seconds to complete
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	if err := srv.Shutdown(ctx2); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
