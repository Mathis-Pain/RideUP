package handlers

import (
	"RideUP/models"
	"RideUP/sessions"
	"RideUP/utils"
	"RideUP/utils/getdata"
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

var RideUpAdminHtml = template.Must(template.ParseFiles(
	"templates/rideupadmin.html",
	"templates/inithtml/inithead.html",
	"templates/inithtml/initnav.html",
	"templates/inithtml/initfooter.html",
))

func RideUpAdminHandler(w http.ResponseWriter, r *http.Request) {
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

	isAdmin, err := getdata.IsUserAdmin(db, userID)
	if err != nil || !isAdmin {
		// Redirige vers la page d'accueil si pas admin
		http.Redirect(w, r, "/RideUp", http.StatusSeeOther)
		return
	}
	allRows, err := db.Query(`
	SELECT e.id, e.title, e.description, e.created_by, u.username, e.created_at,
	       e.latitude, e.longitude, e.address, e.start_datetime, e.end_datetime, e.participants
	FROM events e
	JOIN users u ON e.created_by = u.id
	ORDER BY e.start_datetime ASC`)
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
			&e.CreatorName,
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
		availableEvents = append(availableEvents, e)
	}
	var dashBoard models.AdminDashBoard

	// Compter le nombre total d’utilisateurs
	err = db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&dashBoard.NumberUsers)
	if err != nil {
		log.Println("Erreur COUNT users:", err)
		return
	}

	// Compter le nombre d’utilisateurs connectés
	err = db.QueryRow(`SELECT COUNT(*) FROM sessions `).Scan(&dashBoard.NumberUserConnected)
	if err != nil {
		log.Println("Erreur COUNT connected users:", err)
		return
	}

	// Compter le nombre d’événements
	err = db.QueryRow(`SELECT COUNT(*) FROM events`).Scan(&dashBoard.NumberEvent)
	if err != nil {
		log.Println("Erreur COUNT events:", err)
		return
	}

	data := struct {
		AvailableEvents []models.Event
		DashBoard       models.AdminDashBoard
	}{
		AvailableEvents: availableEvents,
		DashBoard:       dashBoard,
	}

	err = RideUpAdminHtml.Execute(w, data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template RideUpAdmin: %v", err)
		utils.NotFoundHandler(w)
	}
}
