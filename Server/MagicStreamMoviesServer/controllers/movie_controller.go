package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/Kudwa-Abhishek/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)




// creating instance/object of validator
var validate = validator.New()


// this func used to return collection of movie data queried from mongodb database and return to client
// Capital func name as we need this func to be public.
func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.JSON(200, gin.H{"message"   -> SImply check on localhost.
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// creating connection to our movies collection
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

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
func GetMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		//Read in
		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required!"})
			return
		}
		var movie models.Movie

		// creating connection to our movies collection
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)


		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		c.JSON(http.StatusOK, movie)
	}
}

// Adding movie data from client to movies collection. POST using REST API.
func AddMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
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

		// creating connection to our movies collection
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		//code for inserting data into Mongodb collection passed in from client
		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}

// Updating admin review for a movie
func AdminReviewUpdate(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role not found in context"})
			return
		}
		if role != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user must be part of the admin role"})
			return
		}

		movieId := c.Param("imdb_id") // getting movie id from url param
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie Id is required"})
			return
		}
		//local struct named req for request
		var req struct {
			AdminReview string `json:"admin_review"`
		}
		//local struct for response
		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}
		//binding json data(hhtp) sent by client to req struct
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		// Using GetReviewRanking() to extract sentiment from admin review usig langchain go package as done in below func
		sentiment, rankVal, err := GetReviewRanking(req.AdminReview, client, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting review ranking", "details": err.Error()})
			return
		}
		//updating the movie document in Mongodb with new admin review and ranking
		filter := bson.M{"imdb_id": movieId}
		//actual update which tells mongodb what fields to update
		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_name":  sentiment,
					"ranking_value": rankVal,
				},
			},
		}
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// creating connection to our movies collection
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)


		//Update one method to update single document in Mongodb
		result, err := movieCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie"})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		//preparing response to send back to client
		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, resp)
	}
}

// func that will encapsulate the functionality for prompting the AI and getting the ranking
func GetReviewRanking(admin_review string, client *mongo.Client, c *gin.Context) (string, int, error) {
	rankings, err := GetRankings(client, c)

	if err != nil {
		return "", 0, err
	}
	//creating a variable
	sentimentDelimited := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited += ranking.RankingName + ", "
		}
	}
	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	err = godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}
	//communicate with OpenAI using the unique API key saved in the .env file
	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")

	if OpenAiApiKey == "" {
		return "", 0, errors.New("could not read OPENAI_API_KEY")
	}

	llm, err := openai.New(openai.WithToken(OpenAiApiKey))

	if err != nil {
		return "", 0, err
	}
	//read base prompt template from .env file
	base_prompt_template := os.Getenv("BASE_PROMPT_TEMPLATE")

	//we want to replace {rankings} in .env with actual rankings from mongodb
	base_prompt := strings.Replace(base_prompt_template, "{rankings}", sentimentDelimited, 1)
	//actually calling OpenAi to get response
	response, err := llm.Call(context.Background(), base_prompt+admin_review)

	if err != nil {
		return "", 0, err
	}
	rankVal := 0

	//verification with for loop to check that the ranking sent back from prompt is indeed in rankings collection
	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}
	return response, rankVal, nil
}

// func getting all the rankings from mongodb and placing it in prompt in .env
func GetRankings(client *mongo.Client, c *gin.Context) ([]models.Ranking, error) {
	var rankings []models.Ranking

	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	// creating ranking collection
	var rankingCollection *mongo.Collection = database.OpenCollection("rankings", client)


	cursor, err := rankingCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	//putting whatevs in the cursor into rankings slice at a particular address
	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}
	return rankings, nil
}
func GetRecommendedMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		//extract user_id value from context
		userId, err := utils.GetUserIdFromContext(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User id not found in context"})
			return
		}

		//get user favourite genres
		favourite_genres, err := GetUsersFavouriteGenres(userId, client, c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Warning: .env not found")
		}
		var recommendedMovieLimitVal int64 = 5 //only want 5

		recommendedMovieLimitStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")

		if recommendedMovieLimitStr != "" {
			recommendedMovieLimitVal, _ = strconv.ParseInt(recommendedMovieLimitStr, 10, 64)
		}

		findOptions := options.Find()

		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})

		findOptions.SetLimit(recommendedMovieLimitVal)

		filter := bson.M{"genre.genre_name": bson.M{"$in": favourite_genres}}

		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// creating connection to our movies collection
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)


		cursor, err := movieCollection.Find(ctx, filter, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
			return
		}
		defer cursor.Close(ctx)

		var recommendedMovies []models.Movie

		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, recommendedMovies)
	}
}

// get user favourite genres
func GetUsersFavouriteGenres(userId string, client *mongo.Client, c *gin.Context) ([]string, error) {
	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}

	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}

	opts := options.FindOne().SetProjection(projection)

	var result bson.M

	var userCollection *mongo.Collection = database.OpenCollection("users", client)

	err := userCollection.FindOne(ctx, filter, opts).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
	}
	favGenresArray, ok := result["favourite_genres"].(bson.A)
	if !ok {
		return []string{}, errors.New("unable to retrieve favourite genres for user")
	}

	var genreNames []string

	for _, item := range favGenresArray {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}
	return genreNames, nil
}
func GetGenres(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var genreCollection *mongo.Collection = database.OpenCollection("genres", client)

		cursor, err := genreCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching movie genres"})
			return
		}
		defer cursor.Close(ctx)

		var genres []models.Genre
		if err := cursor.All(ctx, &genres); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, genres)

	}
}
