package handlers

import (
	"RideUP/utils"
	"html/template"
	"log"
	"net/http"
)

var NewRideHtml = template.Must(template.ParseFiles("templates/newride.html", "templates/inithtml/inithead.html", "templates/inithtml/initnav.html", "templates/inithtml/initfooter.html"))

func NewRideHandler(w http.ResponseWriter, r *http.Request) {
	err := NewRideHtml.Execute(w, nil)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template NewRideHtml: %v", err)
		utils.NotFoundHandler(w)
	}
}
