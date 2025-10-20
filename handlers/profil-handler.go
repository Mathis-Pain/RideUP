package handlers

import (
	"RideUP/utils"
	"html/template"
	"log"
	"net/http"
)

var ProfilHtml = template.Must(template.ParseFiles(
	"templates/profil.html",
	"templates/inithtml/inithead.html",
	"templates/inithtml/initnav.html",
	"templates/inithtml/initfooter.html",
))

func ProfilHandler(w http.ResponseWriter, r *http.Request) {
	if err := ProfilHtml.Execute(w, nil); err != nil {
		log.Printf("Erreur lors de l'exécution du template rideup.html: %v", err)
		utils.InternalServError(w)
	}
}
