package models

type MongoUser struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	RestID   int    `json:"restid" bson:"restid"`
}
