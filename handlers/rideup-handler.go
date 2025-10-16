package handlers

import (
	"RideUP/models"
	"RideUP/sessions"
	"RideUP/utils"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var RideUpHtml = template.Must(template.ParseFiles(
	"templates/rideup.html",
	"templates/inithtml/inithead.html",
	"templates/inithtml/initnav.html",
	"templates/inithtml/initfooter.html",
))

func RideUpHandler(w http.ResponseWriter, r *http.Request) {

	session, err := sessions.GetSessionFromRequest(r)
	if err != nil {
		log.Printf("Erreur : pas d'utilisateur connecté")
		http.Redirect(w, r, "/Connect", http.StatusSeeOther)
		return
	}

	userID := session.UserID

	db, err := sql.Open("sqlite3", "./data/RideUp.db")
	if err != nil {
		utils.InternalServError(w)
		return
	}
	defer db.Close()

	// -----------------------------
	// 🔹 Sorties créées par l'utilisateur
	// -----------------------------
	rows, err := db.Query("SELECT id, title, description, created_by, created_at, latitude, longitude, start_datetime, end_datetime, max_participants FROM events WHERE created_by = ?", userID)
	if err != nil {
		log.Println("Erreur SELECT userEvents:", err)
	}
	defer rows.Close()

	var userEvents []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.CreatedBy,
			&e.CreatedAt,
			&e.Latitude,
			&e.Longitude,
			&e.StartDatetime,
			&e.EndDatetime,
			&e.MaxParticipants,
		); err != nil {
			log.Println("Erreur Scan userEvents:", err)
			continue
		}
		userEvents = append(userEvents, e)
	}

	// -----------------------------
	// 🔹 Sorties créées par d'autres utilisateurs
	// -----------------------------
	rows2, err := db.Query("SELECT id, title, description, created_by, created_at, latitude, longitude, start_datetime, end_datetime, max_participants FROM events WHERE created_by != ?", userID)
	if err != nil {
		log.Println("Erreur SELECT availableEvents:", err)
	}
	defer rows2.Close()

	var availableEvents []models.Event
	for rows2.Next() {
		var e models.Event
		if err := rows2.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.CreatedBy,
			&e.CreatedAt,
			&e.Latitude,
			&e.Longitude,
			&e.StartDatetime,
			&e.EndDatetime,
			&e.MaxParticipants,
		); err != nil {
			log.Println("Erreur Scan availableEvents:", err)
			continue
		}
		availableEvents = append(availableEvents, e)
	}

	// -----------------------------
	// 🔹 Données envoyées au template
	// -----------------------------
	data := struct {
		UserEvents      []models.Event
		AvailableEvents []models.Event
	}{
		UserEvents:      userEvents,
		AvailableEvents: availableEvents,
	}

	err = RideUpHtml.Execute(w, data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template rideup.html: %v", err)
		utils.NotFoundHandler(w)
	}
}
