// create movie struct within this file
// contains various structs representing our data models
package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Defining genre struct
type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id" validate:"required"`                     //unique identify each genre within db
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required, min=2, max=100"` // actual name of genre as string ex: ROMCOM

}

// We can have rules for each of our fields so that valid data is entered. -> ex imdb_id, title... it is done using the keyword validate
// struct for Ranking
type Ranking struct {
	RankingValue int    `bson:"ranking_value" json:"ranking_value" validate:"required"` //assosciated value ex: Good has value 2
	RankingName  string `bson:"ranking_name" json:"ranking_name" validate:"required"`   // can be excellent, okay, good, bad, terrible
}

type Movie struct {
	ID         bson.ObjectID `bson:"_id" json:"_id"`                             //uniquely identifying movie/ doc
	ImdbID     string        `bson:"imdb_id" json:"imdb_id" validate:"required"` //unique ID for episode/person/move
	Title      string        `bson:"title" json:"title" validate:"required, min=2, max=500"`
	PosterPath string        `bson:"poster_path" json:"poster_path" validate:"required, url"` //used to store URL that points to movie poster
	YoutubeID  string        `bson:"youtube_id" json:"youtube_id" validate:"required"`        // points to utube ID / URL which points to appr trailer
	Genre      []Genre       `bson:"genre" json:"genre" validate:"required, dive"`            // array which can haave 2/more values point to single movie
	//dive ensures that the value in Genre struct/ nested structs are provided compulsorily.
	AdminReview string  `bson:"admin_review" json:"admin_review" validate:"required"`
	Ranking     Ranking `bson:"ranking" json:"ranking" validate:"required"`
}
