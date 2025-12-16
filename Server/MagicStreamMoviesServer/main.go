package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	//root path as in go.mod
	"github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	//fmt.Println("Hello go world!")
	// Basic func to test GIn-GONIC installed
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello , MagicStreamMovies!")
	})

	//relevant cod eto check localhost running ig
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file ")
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
			log.Println("Allowed Origin:", origins[i])
		}
	} else {
		origins = []string{"http://localhost:9845"}
		log.Println("Allowed Origin: http://localhost:9845")
	}
	config := cors.Config{}
	config.AllowOrigins = origins
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.Use(gin.Logger())

	var client *mongo.Client = database.Connect()
	
	// good to ping database before we setup routes here.
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("FAILED to reach server: %v", err)
	}

	// make sure we disconnect from the datbase
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("Failed to disconnect from MongoDB : %v", err)
		}
	}()
	//setting up unprotected routes
	routes.SetupUnProtectedRoutes(router, client)
	//setting up protected routes
	routes.SetupProtectedRoutes(router, client)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}

	//We want to create code for facilitation of separating our unprotexted roots from protected routes.(authn)
	//For this we create a file in routes folder called protected_routes.go
}
