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
			} else {
				// En cas d'erreur qui ne vient pas de la base de données
				// Création d'une session temporaire et anonyme
				err := InitSession(w, 0, "LoginErr", loginErr.Error())
				if err != nil {
					utils.InternalServError(w)
					return
				}

				// Redirection vers la page d'origine
				http.Redirect(w, r, referer, http.StatusSeeOther)
				return
			}
		}
		//Invalider toutes les sessions existantes
		if err := sessions.InvalidateUserSessions(user.ID); err != nil {
			utils.InternalServError(w)
			return
		}

		err = InitSession(w, user.ID, "user", user.Username)
		if err != nil {
			utils.InternalServError(w)
			return
		}

		http.Redirect(w, r, referer, http.StatusSeeOther)

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
