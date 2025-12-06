package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// creating connection to our movies collection
var movieCollection *mongo.Collection = database.OpenCollection("movies")

// creating instance/object of validator
var validate = validator.New()

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

// getting a single movie by its id or imdb_id
func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//Read in
		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required!"})
			return
		}
		var movie models.Movie

		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		c.JSON(http.StatusOK, movie)
	}
}

// Adding movie data from client to movies collection. POST using REST API.
func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel() //ensures that when your function timeouts, resources are released.

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		//Validation code after installing go - playground validator package in go.mod
		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		//code for inserting data into Mongodb collection passed in from client
		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}
