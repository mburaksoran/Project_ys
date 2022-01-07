package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/mburaksoran/models"
	"github.com/mburaksoran/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", shared.Config.POSTGRESURL)
	if err != nil {
		log.Fatal(err)
	}

}

func GetOrdersWithRestID(id int) []*models.OrderFromDB {
	rows, err := db.Query("SELECT * FROM orders WHERE rest_id =" + strconv.Itoa(id))
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No Records Found")
			return nil
		}
	}
	defer rows.Close()
	var orders []*models.OrderFromDB
	for rows.Next() {
		prd := &models.OrderFromDB{}
		err := rows.Scan(&prd.OrderID, &prd.OrderTime, &prd.OrderElemCount, &prd.OrderValue, &prd.RestID, &prd.UserID, &prd.ItemID)
		if err != nil {
			log.Fatal(err)
		}
		orders = append(orders, prd)

	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return orders
}

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

func BulkInsert(client *mongo.Client, ctx *context.Context, person models.MongoUser) error {
	collection := client.Database("YS_proje").Collection("YS_Proje_Users")
	_, err := collection.InsertOne(*ctx, person)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func GetDataFromMdb(client *mongo.Client, ctx *context.Context, username string) models.MongoUser {
	//var test models.MongoUser
	usersCollection := client.Database("YS_proje").Collection("YS_Proje_Users")

	var user models.MongoUser
	usersCollection.FindOne(*ctx, bson.M{"username": username}).Decode(&user)
	return user
}
