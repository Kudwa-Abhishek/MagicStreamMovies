package routes

import (
	controller "github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// create a func to setup protected routes
func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	router.Use(middleware.AuthMiddleware()) // applying auth middleware to all routes defined in this func

	//endpoints that need protection go here
	router.POST("/addmovie", controller.AddMovie(client))
	router.GET("/movie/:imdb_id", controller.GetMovie(client))
	router.GET("/recommendedmovies", controller.GetRecommendedMovies(client))
	router.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate(client))

}
