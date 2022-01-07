package routers

import (
	"log"
	"net/http"

	"github.com/mburaksoran/handlers"
	"github.com/mburaksoran/helpers"
)

func HandleRequests() {
	http.HandleFunc("/login", handlers.Login)
	http.Handle("/welcome", helpers.CheckIfAuthorized(handlers.WelcomePage))
	http.HandleFunc("/getdata", handlers.GetData)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
