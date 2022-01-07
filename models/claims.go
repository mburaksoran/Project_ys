package models

import jwt "github.com/dgrijalva/jwt-go"

type Claims struct {
	Username   string `json:"username"`
	UserRestID int    `json:"restID"`
	jwt.StandardClaims
}
