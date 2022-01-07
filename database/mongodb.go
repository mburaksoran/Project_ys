package database

import (
	"context"
	"log"
	"time"

	"github.com/mburaksoran/models"
	"github.com/mburaksoran/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Creating client and context for mongo db return values will be used to mongodb operations.
func MongoOpen() (*mongo.Client, *context.Context, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(shared.Config.MONGOURL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return client, &ctx, err
}

// using to insert user information to mongodb.
func InsertDataToMdb(client *mongo.Client, ctx *context.Context, person models.MongoUser) error {
	collection := client.Database("YS_proje").Collection("YS_Proje_Users")
	_, err := collection.InsertOne(*ctx, person)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

//collecting user information from mongodb to compare login information.
func GetDataFromMdb(client *mongo.Client, ctx *context.Context, username string) models.MongoUser {

	usersCollection := client.Database("YS_proje").Collection("YS_Proje_Users")

	var user models.MongoUser
	usersCollection.FindOne(*ctx, bson.M{"username": username}).Decode(&user)
	return user
}
