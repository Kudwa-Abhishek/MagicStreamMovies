package routes

import (
	controller "github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/middleware"
	"github.com/gin-gonic/gin"
)
// create a func to setup protected routes
func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware()) // applying auth middleware to all routes defined in this func

	//endpoints that need protection go here
	router.POST("/addmovie", controller.AddMovie())
	router.GET("/movie/:imdb_id", controller.GetMovie())
}
