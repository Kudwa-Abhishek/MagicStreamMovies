package main

import (
	"fmt"

	//root path as in go.mod
	"github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	//fmt.Println("Hello go world!")
	// Basic func to test GIn-GONIC installed
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello , MagicStreamMovies!")
	})
	//setting up unprotected routes
	routes.SetupUnProtectedRoutes(router)
	//setting up protected routes
	routes.SetupProtectedRoutes(router)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}

	//We want to create code for facilitation of separating our unprotexted roots from protected routes.(authn)
	//For this we create a file in routes folder called protected_routes.go
}
