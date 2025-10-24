package handlers

import (
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

	err = RideUpAdminHtml.Execute(w, nil)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template RideUpAdmin: %v", err)
		utils.NotFoundHandler(w)
	}
}
