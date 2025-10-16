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

	_ "github.com/mattn/go-sqlite3"
)

var RideUpHtml = template.Must(template.ParseFiles(
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
	// 🔹 Sorties créées par l'utilisateur
	// -----------------------------
	userRows, err := db.Query(`
		SELECT id, title, description, created_by, created_at, latitude, longitude, start_datetime, end_datetime, max_participants
		FROM events
		WHERE created_by = ?`, userID)
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
	// 🔹 Toutes les sorties disponibles
	// -----------------------------
	allRows, err := db.Query(`
		SELECT id, title, description, created_by, created_at, latitude, longitude, start_datetime, end_datetime, max_participants
		FROM events`)
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
	// 🔹 Conversion latitude/longitude → adresse
	// -----------------------------
	for i := range userEvents {
		address, err := utils.ReverseGeocodeSimple(userEvents[i].Latitude, userEvents[i].Longitude)
		if err != nil || address == nil {
			// Si erreur, créer une adresse avec les coordonnées
			userEvents[i].Location = &models.SimpleAddress{
				Rue: fmt.Sprintf("Coordonnées: %.6f, %.6f", userEvents[i].Latitude, userEvents[i].Longitude),
			}
		} else {
			// Assigner directement le pointeur
			userEvents[i].Location = address
		}
	}

	for i := range availableEvents {
		address, err := utils.ReverseGeocodeSimple(availableEvents[i].Latitude, availableEvents[i].Longitude)
		if err != nil || address == nil {
			// Si erreur, créer une adresse avec les coordonnées
			availableEvents[i].Location = &models.SimpleAddress{
				Rue: fmt.Sprintf("Coordonnées: %.6f, %.6f", availableEvents[i].Latitude, availableEvents[i].Longitude),
			}
		} else {
			// Assigner directement le pointeur
			availableEvents[i].Location = address
		}
	}

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
		AvailableEvents: availableEvents,
	}

	if err := RideUpHtml.Execute(w, data); err != nil {
		log.Printf("Erreur lors de l'exécution du template rideup.html: %v", err)
		utils.InternalServError(w)
	}
}
