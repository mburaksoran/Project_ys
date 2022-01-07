package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mburaksoran/database"
	"github.com/mburaksoran/models"
	"github.com/mburaksoran/utils"
)

var jwtKey = []byte("secret_key")
var SendRecords []models.DataForClient
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	var html models.HtmlInformation
	if r.Method == http.MethodPost {
		uname := r.FormValue("username")
		pwd := r.FormValue("password")
		credentials.Username = uname
		credentials.Password = utils.PassToSHA256(pwd)
		client, ctx, err := database.MongoOpen()
		if err != nil {
			html.Error = "Lütfen daha sonra tekrar deneyin"
			tpl.ExecuteTemplate(w, "login.gohtml", html)
		}
		user := database.GetDataFromMdb(client, ctx, uname)

		if credentials.Password != user.Password || credentials.Username != user.Username {
			html.Error = "Kullanıcı adı veya şifre hatalı"
			tpl.ExecuteTemplate(w, "login.gohtml", html)
			return
		}

		expirationTime := time.Now().Add(time.Minute * 5)

		claims := &models.Claims{
			Username:   credentials.Username,
			UserRestID: user.RestID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedTokenString, err := token.SignedString(jwtKey)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// test değişkkeninin ismi dğeişecek
		test := http.Cookie{
			Name:    "token",
			Value:   signedTokenString,
			Expires: expirationTime,
		}
		http.SetCookie(w, &test)

		http.Redirect(w, r, "/welcome", http.StatusSeeOther)

	}

	tpl.ExecuteTemplate(w, "login.gohtml", nil)

}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
		}

		tokenstr := cookie.Value

		claims := &models.Claims{}

		tkn, err := jwt.ParseWithClaims(tokenstr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
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

		if tkn.Valid {
			endpoint(w, r)
		}

	})
}

func WelcomePage(w http.ResponseWriter, r *http.Request) {
	var html models.HtmlInformation
	if r.Method == http.MethodPost {
		restid, username := utils.GetRestIDFromToken(r, jwtKey)
		filepath := utils.UploadData(r)
		dataFromClient := utils.UploadedDataToStruct(filepath)
		datafromdatabase := database.GetOrdersWithRestID(restid)
		SendRecords = utils.DataPrepForClient(dataFromClient, datafromdatabase)

		html.Username = username
		html.Status = "Dosya başarıyla yüklendi"
		html.Error = ""
		tpl.ExecuteTemplate(w, "welcome.gohtml", html)

	} else {
		_, username := utils.GetRestIDFromToken(r, jwtKey)
		html.Username = username
		html.Status = ""
		tpl.ExecuteTemplate(w, "welcome.gohtml", html)
	}

}

func GetData(w http.ResponseWriter, r *http.Request) {
	//header get olarak ayarlanacak
	w.Header().Set("Content-Type", "application/json")
	users := SendRecords
	json.NewEncoder(w).Encode(users)
	SendRecords = make([]models.DataForClient, 0)
}

func HandleRequests() {
	http.HandleFunc("/login", Login)
	http.Handle("/welcome", isAuthorized(WelcomePage))
	http.HandleFunc("/getdata", GetData)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
