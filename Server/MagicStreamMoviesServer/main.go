package main

import (
	"fmt"

	controller "github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers" //root path as in go.mod
	"github.com/gin-gonic/gin"
)

func main() {
	//fmt.Println("Hello go world!")
	// Basic func to test GIn-GONIC installed
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello , MagicStreamMovies!")
	})

	//New endpoint
	router.GET("/movies", controller.GetMovies())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
