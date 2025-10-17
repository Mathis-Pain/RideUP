package handlers

import (
	"RideUP/models"
	"RideUP/sessions"
	"RideUP/utils"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var NewEventHtml = template.Must(template.ParseFiles("templates/newevent.html", "templates/inithtml/inithead.html", "templates/inithtml/initnav.html", "templates/inithtml/initfooter.html"))

func NewEventHandler(w http.ResponseWriter, r *http.Request) {

	// Récupère l'utilisateur connecté
	session, err := sessions.GetSessionFromRequest(r)
	if err != nil {
		log.Printf("Erreur pas d'utilisateur connecté")
		//  http.StatusSeeOther Renvoi un code 303 a l'entete et permet de ne demander un get pour recuperer l'url ce qui  evite de renvoyer un formulaire par exemple
		http.Redirect(w, r, "/Connect", http.StatusSeeOther)
		return
	}

	db, err := sql.Open("sqlite3", "./data/RideUp.db")
	if err != nil {
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	var data models.MapData
	err = db.QueryRow(`SELECT latitude,longitude FROM users WHERE id =?`, session.UserID).Scan(&data.Latitude, &data.Longitude)
	if err != nil {
		log.Println("Erreur récupération coordonnées:", err)
		// Valeurs par défaut (ex: Paris)
		data.Latitude = 48.8566
		data.Longitude = 2.3522
	}

	if r.Method == "POST" {
		// Récupération des valeurs du formulaire
		if err := r.ParseForm(); err != nil {
			log.Printf("ERREUR : ParseForm: %v", err)
			http.Error(w, "Erreur lors de la lecture du formulaire", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		dateStr := r.FormValue("date")
		timeStr := r.FormValue("time")
		location := r.FormValue("location")

		if title == "" || dateStr == "" || timeStr == "" || location == "" {
			http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
			return
		}
		// Conversion date + heure en time.Time
		startStr := dateStr + " " + timeStr
		startDatetime, err := time.Parse("2006-01-02 15:04", startStr)
		if err != nil {
			log.Printf("ERREUR : parsing datetime: %v", err)
			http.Error(w, "Date ou heure invalide", http.StatusBadRequest)
			return
		}
		// Extraction latitude et longitude du champ location
		var lat, lon float64
		_, err = fmt.Sscanf(location, "Lat: %f, Lon: %f", &lat, &lon)
		if err != nil {
			log.Printf("ERREUR : parsing location: %v", err)
			http.Redirect(w, r, "/NewEvent", http.StatusSeeOther)
			return
		}
		session, err := sessions.GetSessionFromRequest(r)
		if err != nil {
			log.Printf("Erreur pas d'utilisateur connecté")
			http.Redirect(w, r, "/Connect", http.StatusSeeOther)
			return
		}
		// Insertion dans la base de données
		_, err = db.Exec(`
        INSERT INTO events (title, created_by, latitude, longitude, start_datetime)
        VALUES (?, ?, ?, ?, ?)
    `, title, session.UserID, lat, lon, startDatetime)
		if err != nil {
			log.Printf("ERREUR : insertion event: %v", err)
			http.Error(w, "Impossible de créer l'événement", http.StatusInternalServerError)
			return
		}
		// Redirection vers la page des événements ou confirmation
		http.Redirect(w, r, "/RideUp", http.StatusSeeOther)
	}
	err = NewEventHtml.Execute(w, data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template NewRideHtml: %v", err)
		utils.NotFoundHandler(w)
	}
}
