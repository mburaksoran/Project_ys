package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mburaksoran/database"
	"github.com/mburaksoran/helpers"
	"github.com/mburaksoran/models"
	"github.com/mburaksoran/utils"
)

// To sign the token, the function needs to secret key, for using multiple places, secret key arrange as a global variable in this case.
//it could be used in .env as a SECREC_KEY tag but i chosed to use like that.
var jwtKey = []byte("secret_key")

// in the project, there are two different type of data used. to combine them and using in seperate functions, it arrange as a global variable .
var SendRecords []models.DataForClient
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

//in login function, firstly creating two variable. variable that named credentials used for the data which came from html form and pased a struct, html is used for informing the client.
//for example if there is an error, the client inform with html.status.
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	var html models.HtmlInformation
	if r.Method == http.MethodPost {
		// gathering username and password from html form and arrange them to credentials, in the project all passwords hashed with sha256.
		credentials.Username = r.FormValue("username")
		credentials.Password = helpers.PassToSHA256(r.FormValue("password"))
		client, ctx, err := database.MongoOpen()
		// contolling the connection to the mongodb if there is an error with connection, client will inform.
		if err != nil {
			html.Status = "Bağlantı sırasında bir hata meydana geldi, lütfen daha sonra tekrar deneyin"
			tpl.ExecuteTemplate(w, "login.gohtml", html)
		}
		// after the successful connection, username and password will be compare with data which came from mongodb. if usernames and passwords are not same, client will be inform about
		// wrong username or password .
		user := database.GetDataFromMdb(client, ctx, credentials.Username)
		ok, status := helpers.LoginCheck(credentials, user)
		if !ok {
			html.Status = status
			tpl.ExecuteTemplate(w, "login.gohtml", html)
			return
		}
		// creating expiration time and claims for authentication token. it arrange 5 minutes expirationtime.
		expirationTime := time.Now().Add(time.Minute * 5)
		claims := &models.Claims{
			Username:   credentials.Username,
			UserRestID: user.RestID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		// sign the token with secret key via jwt package. if there is an error with signing client will be inform.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedTokenString, err := token.SignedString(jwtKey)
		if err != nil {
			html.Status = "Giriş Sırasında bir hata meydana geldi. Lütfen tekrar giriş yapınız"
			tpl.ExecuteTemplate(w, "login.gohtml", html)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:    "token",
			Value:   signedTokenString,
			Expires: expirationTime,
		}
		// setting authentication token as a cookie.
		http.SetCookie(w, &cookie)
		// after setting the cookie, the client will redirect to welcome page.
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)

	}
	// to filling the login information, blank html form will send to client. after filling the username and password and hit the submit button
	// the html form will send as post-request to server.
	tpl.ExecuteTemplate(w, "login.gohtml", nil)

}

// after successful login, client will be redirect to the welcome page. in welcome page , user information will gather from token,
//after that user's restaurant id will use for connecting and collecting
// to the database. the data which gather from postgres-sql parse in a struct via  GetOrdersWithRestID function and
//the data which came from client will be read and save to the temp files as csv file via UploadData function

// with UploadedDataToStruct function, csv data will be parse to struct array.
//both data which are collected from database and client , will be parse to same struct "models.DataForClient" with DataPrepForClient function.
func WelcomePage(w http.ResponseWriter, r *http.Request) {
	var html models.HtmlInformation
	if r.Method == http.MethodPost {
		restid, username := helpers.GetInformationFromToken(r, jwtKey)
		filepath, err := utils.UploadData(r)
		//controlling the data which sended from client. if there is an error in uploading phase , client will be inform.
		if err != nil {
			html.Username = username
			html.Status = "Dosya yükleme sırasında bir hata meydana geldi"
			tpl.ExecuteTemplate(w, "welcome.gohtml", html)
			return
		}
		//controlling the uploaded file's type if data type wrong client will inform about data type.
		dataFromClient, err := utils.UploadedDataToStruct(filepath)
		if err != nil {
			html.Username = username
			html.Status = "Yüklenen dosya uygun formatta değil"
			tpl.ExecuteTemplate(w, "welcome.gohtml", html)
			return
		}
		// gathering Restaurant's data from database. if there is a connection error client will be inform.
		datafromdatabase, err := database.GetOrdersWithRestID(restid)
		if err != nil {
			html.Username = username
			html.Status = "Veritabanına erişilemiyor lütfen daha sonra tekrar deneyin"
			tpl.ExecuteTemplate(w, "welcome.gohtml", html)
			return
		}
		// creating combined data which came from database and client.
		SendRecords = utils.DataPrepForClient(dataFromClient, datafromdatabase)
		html.Username = username
		html.Status = "Dosya başarıyla yüklendi"
		tpl.ExecuteTemplate(w, "welcome.gohtml", html)

	} else {
		_, username := helpers.GetInformationFromToken(r, jwtKey)
		html.Username = username
		html.Status = ""
		tpl.ExecuteTemplate(w, "welcome.gohtml", html)
	}

}

// prepared data will be serve with Getdata method. when client successfully upload their data, api show a button that named see the result.
// when client click the button it will be pop up to R shiny services.
//when the R shiny services triggered, it will send a get request to the golang api and gather combined data from go api .
func GetData(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		users := SendRecords
		json.NewEncoder(w).Encode(users)
		SendRecords = make([]models.DataForClient, 0)
	}

}
