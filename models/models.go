package models

import jwt "github.com/dgrijalva/jwt-go"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username   string `json:"username"`
	UserRestID int    `json:"restID"`
	jwt.StandardClaims
}

type HtmlInformation struct {
	Username string
	Status   string
	Error    string
}

type OrderFromDB struct {
	OrderID        int
	OrderTime      string
	OrderElemCount int
	OrderValue     int
	RestID         int
	UserID         int
	ItemID         []uint8
}

type OrderFromClient struct {
	OrderID   int
	SentBy    string
	OrderTime string
	ItemVal   float64
	ClosedBy  string
	RestName  string
	PaidType  string
	From      string
	Item      string
	ItemCount int
	TableName string
	ClosedAt  string
}

type DataForClient struct {
	OrderTime    string  `json:"time"`
	TotalPayment float64 `json:"payment"`
	DataComeFrom int     `json:"datacomefrom"`
}

type MongoUser struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	RestID   int    `json:"restid" bson:"restid"`
}
