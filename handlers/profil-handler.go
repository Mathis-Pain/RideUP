package handlers

import (
	"RideUP/models"
	"RideUP/sessions"
	"RideUP/utils"
	"database/sql"
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

	// Récupère l'utilisateur connecté
	session, err := sessions.GetSessionFromRequest(r)
	if err != nil {
		log.Printf("Erreur pas d'utilisateur connecté")
		http.Redirect(w, r, "/Connect", http.StatusSeeOther)
		return
	}

	db, err := sql.Open("sqlite3", "./data/RideUp.db")
	if err != nil {
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	// Coordonnées par défaut ou de l'utilisateur
	var data models.MapData
	err = db.QueryRow(`SELECT latitude, longitude FROM users WHERE id = ?`, session.UserID).
		Scan(&data.Latitude, &data.Longitude)
	if err != nil {
		log.Println("Erreur récupération coordonnées:", err)
		data.Latitude = 48.8566
		data.Longitude = 2.3522
	}

	if err := ProfilHtml.Execute(w, data); err != nil {
		log.Printf("Erreur lors de l'exécution du template rideup.html: %v", err)
		utils.InternalServError(w)
	}
}
