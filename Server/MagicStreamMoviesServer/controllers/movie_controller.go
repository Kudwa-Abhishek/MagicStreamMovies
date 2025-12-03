package controllers

import (
	"github.com/gin-gonic/gin"
)

// this func used to return collection of movie data queried from mongodb database and return to client
// Capital func name as we need this func to be public.
func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "List of movies"})
	}
}
