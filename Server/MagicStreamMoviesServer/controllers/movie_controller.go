package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// creating connection to our movies collection
var movieCollection *mongo.Collection = database.OpenCollection("movies")

// this func used to return collection of movie data queried from mongodb database and return to client
// Capital func name as we need this func to be public.
func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.JSON(200, gin.H{"message"   -> SImply check on localhost.
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies []models.Movie

		//cursor for purpose of querying and returning data from mongodb database
		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies"})
		}

		c.JSON(http.StatusOK, movies)
	}
}

//Running this will give JSON data of movies from Mongodb on endpoint /movies in localhost.
// As of when I am running, it shows the JSON format of movies which will look better once REACT is connected.
