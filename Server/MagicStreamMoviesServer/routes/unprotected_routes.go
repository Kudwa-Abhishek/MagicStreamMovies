package routes

import (
	controller "github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/gin-gonic/gin"
)
// create a func to setup protected routes
func SetupUnProtectedRoutes(router *gin.Engine) {
	//endpoints that do not need protection go here
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())
	router.GET("/movies", controller.GetMovies())
}