package handlers

import (
	"RideUP/utils"
	"html/template"
	"log"
	"net/http"
)

var UpdatePasswordHtml = template.Must(template.ParseFiles(
	"templates/updatepassword.html",
	"templates/inithtml/inithead.html",
	"templates/inithtml/initnav.html",
	"templates/inithtml/initfooter.html",
))

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	err := UpdatePasswordHtml.Execute(w, nil)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template UpdatePasswordHtml: %v", err)
		utils.NotFoundHandler(w)
	}
}
