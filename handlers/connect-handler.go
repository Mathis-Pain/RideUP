package handlers

import (
	"RideUP/sessions"
	"RideUP/utils"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var ConnectionHtml = template.Must(template.ParseFiles("templates/connection.html"))

func ConnectHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	// Gestion de la methode GET
	if r.Method == "GET" {
		if err := ConnectionHtml.Execute(w, nil); err != nil {
			log.Printf("ERREUR : Affichage page connexion: %v", err)
			utils.InternalServError(w)
			return
		}
		return
	}

	// Gestion de la methode Post
	if r.Method == "POST" {
		// Récupération des valeurs du formulaire
		if err := r.ParseForm(); err != nil {
			log.Printf("ERREUR : ParseForm: %v", err)
			http.Error(w, "Erreur lors de la lecture du formulaire", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Ouverture de la DB
		db, err := sql.Open("sqlite3", "./data/RideUp.db")
		if err != nil {
			utils.InternalServError(w)
			log.Printf("Erreur : Ouverture base de données: %v", err)
			return
		}
		defer db.Close()

		user, loginErr := utils.Authentification(db, email, password)
		if loginErr != nil {
			if strings.Contains(loginErr.Error(), "db") {
				// En cas d'erreur dans la base de données
				utils.InternalServError(w)
			}
			// Mauvais identifiant / mot de passe → on ne crée pas de vraie session
			http.Redirect(w, r, referer+"?error="+loginErr.Error(), http.StatusSeeOther)
			return
		}

		// Ici, tout est ok → créer la vraie session
		if err := sessions.InvalidateUserSessions(user.ID); err != nil {
			utils.InternalServError(w)
			return
		}
		if err := InitSession(w, user.ID, "user", user.Email); err != nil {
			utils.InternalServError(w)
			return
		}

		http.Redirect(w, r, "/RideUp", http.StatusSeeOther)
		return
	}
}

// Crée, sauvegarde une session, permet d'y insérer une donnée et pose le cookie
func InitSession(w http.ResponseWriter, id int, fieldName string, fieldData any) error {
	session, err := sessions.CreateSession(id)
	if err != nil {
		return err
	}

	session.Data[fieldName] = fieldData
	if err := sessions.SaveSessionToDB(session); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   false, // false en local, true si HTTPS
		Path:     "/",
	})
	return nil
}
