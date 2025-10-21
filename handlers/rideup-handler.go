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

var EventHtml = template.Must(template.ParseFiles(
	"templates/rideup.html",
	"templates/inithtml/inithead.html",
	"templates/inithtml/initnav.html",
	"templates/inithtml/initfooter.html",
))

func RideUpHandler(w http.ResponseWriter, r *http.Request) {
	// 🔹 Récupération de la session
	session, err := sessions.GetSessionFromRequest(r)
	if err != nil {
		log.Printf("Erreur : pas d'utilisateur connecté")
		http.Redirect(w, r, "/Connect", http.StatusSeeOther)
		return
	}
	userID := session.UserID

	// 🔹 Connexion à la DB
	db, err := sql.Open("sqlite3", "./data/RideUp.db")
	if err != nil {
		utils.InternalServError(w)
		return
	}
	defer db.Close()
	// -----------------------------
	// 🔹suppression des sorties qui sont passées
	// -----------------------------
	_, err = db.Exec(`DELETE FROM events WHERE start_datetime < datetime('now')`)
	if err != nil {
		log.Printf("Erreur suppression événements passés : %v", err)
		utils.InternalServError(w)
		return
	}
	// -----------------------------
	// 🔹 Sorties créées par l'utilisateur
	// -----------------------------
	userRows, err := db.Query(`
	SELECT id, title, description, created_by, created_at, 
	       latitude, longitude, address, start_datetime, end_datetime, participants
	FROM events
	WHERE created_by = ?
	ORDER BY start_datetime ASC`, userID)
	if err != nil {
		log.Println("Erreur SELECT userEvents:", err)
		utils.InternalServError(w)
		return
	}
	defer userRows.Close()

	var userEvents []models.Event
	for userRows.Next() {
		var e models.Event
		if err := userRows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.CreatedBy,
			&e.CreatedAt,
			&e.Latitude,
			&e.Longitude,
			&e.Address,
			&e.StartDatetime,
			&e.EndDatetime,
			&e.Participants,
		); err != nil {
			log.Println("Erreur Scan userEvents:", err)
			continue
		}
		// -----------------------------
		// 🔹 Vérifie si l'utilisateur a rejoint cet event
		// -----------------------------
		var count int
		err = db.QueryRow(`SELECT COUNT(*) FROM event_participants WHERE user_id = ? AND event_id = ?`,
			userID, e.ID).Scan(&count)
		if err == nil && count > 0 {
			e.UserJoined = true
		} else {
			e.UserJoined = false
		}

		userEvents = append(userEvents, e)
	}

	// -----------------------------
	// 🔹 Toutes les sorties disponibles
	// -----------------------------
	allRows, err := db.Query(`
	SELECT id, title, description, created_by, created_at, 
	       latitude, longitude, address, start_datetime, end_datetime, participants
	FROM events
	ORDER BY start_datetime ASC`)
	if err != nil {
		log.Println("Erreur SELECT availableEvents:", err)
		utils.InternalServError(w)
		return
	}
	defer allRows.Close()

	var availableEvents []models.Event
	for allRows.Next() {
		var e models.Event
		if err := allRows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.CreatedBy,
			&e.CreatedAt,
			&e.Latitude,
			&e.Longitude,
			&e.Address,
			&e.StartDatetime,
			&e.EndDatetime,
			&e.Participants,
		); err != nil {
			log.Println("Erreur Scan availableEvents:", err)
			continue
		}
		// -----------------------------
		// 🔹 Vérifie si l'utilisateur a rejoint cet event
		// -----------------------------
		var count int
		err = db.QueryRow(`SELECT COUNT(*) FROM event_participants WHERE user_id = ? AND event_id = ?`,
			userID, e.ID).Scan(&count)
		if err == nil && count > 0 {
			e.UserJoined = true
		} else {
			e.UserJoined = false
		}

		availableEvents = append(availableEvents, e)
	}
	// -----------------------------
	// 🔹 Filtrer les event en fonction des preferences utilisateur
	// -----------------------------
	// Déclaration des variables pour stocker les valeurs
	var latitude, longitude float64
	var preference int

	// Récupération des infos depuis la table users
	err = db.QueryRow(`
    SELECT latitude, longitude, preference 
    FROM users 
    WHERE id = ?`, userID).Scan(&latitude, &longitude, &preference)
	if err != nil {
		log.Printf("Erreur récupération infos utilisateur: %v", err)
		utils.InternalServError(w)
		return
	}

	availableEventsFilter := utils.FilterPreference(availableEvents, latitude, longitude, preference)

	// -----------------------------
	// 🔹 Données envoyées au template
	// -----------------------------
	data := struct {
		ActivePage      string
		UserEvents      []models.Event
		AvailableEvents []models.Event
	}{
		ActivePage:      "RideUp",
		UserEvents:      userEvents,
		AvailableEvents: availableEventsFilter,
	}

	if err := EventHtml.Execute(w, data); err != nil {
		log.Printf("Erreur lors de l'exécution du template rideup.html: %v", err)
		utils.InternalServError(w)
	}
}
