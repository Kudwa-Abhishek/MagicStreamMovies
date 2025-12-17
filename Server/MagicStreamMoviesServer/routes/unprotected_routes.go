package routes

import (
	controller "github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// create a func to setup protected routes
func SetupUnProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	//endpoints that do not need protection go here
	router.POST("/register", controller.RegisterUser(client))
	router.POST("/login", controller.LoginUser(client))
	router.GET("/movies", controller.GetMovies(client))
	router.POST("/logout", controller.LogoutHandler(client))
	router.GET("/genres", controller.GetGenres(client))
	router.POST("/refresh", controller.RefreshTokenHandler(client))
}
