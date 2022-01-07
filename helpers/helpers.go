package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/mburaksoran/models"
)

var jwtKey = []byte("secret_key")

// a function that hashing the string value. for this project it used to hashing passwords that came from client.
func PassToSHA256(passStr string) string {
	// creating new sha256 hasher
	hash := sha256.New()
	hash.Write([]byte(passStr))
	// encrpt the string with hasher.
	sha256Hash := hex.EncodeToString(hash.Sum(nil))

	return sha256Hash
}

// a function that gathering logged client information. Function takes two variable one for gathering request's cookie for user information
// the other is jwtkeys for decrpt signed token which came from cookie,

func GetInformationFromToken(r *http.Request, jwtKey []byte) (int, string) {
	// for parsing information first it needs to claim struct,
	var claim models.Claims
	//  gathering cookie with name "token" from request
	cookie, _ := r.Cookie("token")
	// turning a string cookie's value
	tokenStr := cookie.Value
	// decrpt the signed token and parsing token to claims struct element
	jwt.ParseWithClaims(tokenStr, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	return claim.UserRestID, claim.Username
}

// checking the client information, if username exist and user password is same with password that came from database, it will return true and empty string but if it is not the same
//password, it will return false and information that will inform the client.
func LoginCheck(cred models.Credentials, user models.MongoUser) (bool, string) {

	if cred.Password != user.Password || cred.Username != user.Username {

		return false, "Kullanıcı adı veya şifre hatalı"
	}
	return true, ""
}

// It is a function that controlling authentication token. if it is exist, it will let the client to reach the page if authentication token doesnt exist or expires run off, it will not
// let the client reach the page.
func CheckIfAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// gathering cookies from request.
		cookie, err := r.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
		}
		// assing to cookie's value to tokenstr
		tokenstr := cookie.Value
		// creating a empty claims element to parsing information which will came from tokenstr
		claims := &models.Claims{}
		// parsing the token information to claim if there is an error it will throw an error
		tkn, err := jwt.ParseWithClaims(tokenstr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		// controlling the error if error is not empty it will try to figure out what kind of error occur.
		if err != nil {
			//controlling error if signiature is valid
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// at the end if there is no error and token is valid it will pass to page .
		if tkn.Valid {
			endpoint(w, r)
		}

	})
}
