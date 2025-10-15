package handlers

import (
	"RideUP/utils"
	"html/template"
	"log"
	"net/http"
)

var RideUpHtml = template.Must(template.ParseFiles("templates/rideup.html", "templates/inithtml/inithead.html", "templates/inithtml/initnav.html", "templates/inithtml/initfooter.html"))

func RideUpHandler(w http.ResponseWriter, r *http.Request) {
	err := RideUpHtml.Execute(w, nil)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template rideup.html: %v", err)
		utils.NotFoundHandler(w)
	}
}
