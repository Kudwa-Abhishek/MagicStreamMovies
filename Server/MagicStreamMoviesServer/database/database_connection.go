package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Used to connect gin-gonic web API application to our mongodb database through the use of mongodb for go driver
// Already imported mongo related packages above in import block
func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}
	//read env variable
	MongoDb := os.Getenv("MONGODB_URI")

	if MongoDb == "" {
		log.Fatal("MONGODB_URI not set in .env file")
	}

	fmt.Println("MongoDB URI: ", MongoDb)

	//we shld be able to connect our client code
	clientOptions := options.Client().ApplyURI(MongoDb)

	//actually connect to mongodb database
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil
	}

	return client
}

// create client object based on db instance
var Client *mongo.Client = DBInstance()

// method which opens our actual connection to database.
func OpenCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}
	//read env variable and assign to database
	databaseName := os.Getenv("DATABASE_NAME")

	fmt.Println("DATABASE_NAME:", databaseName)

	//load the collection
	collection := Client.Database(databaseName).Collection(collectionName)
	if collection == nil {
		return nil
	}
	return collection
}
