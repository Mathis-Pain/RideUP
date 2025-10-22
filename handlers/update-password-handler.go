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

	"golang.org/x/crypto/bcrypt"
)

var UpdatePasswordHtml = template.Must(template.ParseFiles(
	"templates/updatepassword.html",
	"templates/inithtml/inithead.html",
	"templates/inithtml/initnav.html",
	"templates/inithtml/initfooter.html",
))

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := UpdatePasswordHtml.Execute(w, nil)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template UpdatePasswordHtml: %v", err)
			utils.NotFoundHandler(w)
		}
		return

	case http.MethodPost:
		session, err := sessions.GetSessionFromRequest(r)
		if err != nil {
			http.Redirect(w, r, "/Connect", http.StatusSeeOther)
			return
		}

		db, err := sql.Open("sqlite3", "./data/RideUp.db")
		if err != nil {
			utils.InternalServError(w)
			return
		}
		defer db.Close()

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Erreur formulaire", http.StatusBadRequest)
			return
		}

		oldPassword := r.FormValue("OldPassword")
		newPassword := r.FormValue("NewPassword")
		confirmPassword := r.FormValue("ConfirmNewPassword")

		// Récupère le hash actuel
		passwordHash, err := getdata.GetPasswordHash(db, session.UserID)
		if err != nil {
			log.Printf("Erreur récupération du mot de passe: %v", err)
			utils.InternalServError(w)
			return
		}

		// Vérifie l'ancien mot de passe
		if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(oldPassword)) != nil {
			data := models.UpdatePassword{
				OldPasswordError: "Ancien mot de passe incorrect",
			}
			UpdatePasswordHtml.Execute(w, data)
			return
		}

		// Vérifie la validité du nouveau mot de passe
		newPassErr := utils.ValidPassword(newPassword, confirmPassword)
		if newPassErr != "" {
			data := models.UpdatePassword{
				NewPasswordError: newPassErr,
			}
			UpdatePasswordHtml.Execute(w, data)
			return
		}

		// Hash le nouveau mot de passe
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			utils.InternalServError(w)
			return
		}

		// Met à jour le mot de passe dans la DB
		_, err = db.Exec(`UPDATE users SET password_hash = ? WHERE id = ?`, hashedPassword, session.UserID)
		if err != nil {
			log.Printf("Erreur de mise à jour du mot de passe: %v", err)
			utils.InternalServError(w)
			return
		}

		// Redirige vers le profil avec succès
		http.Redirect(w, r, "/RideUp", http.StatusSeeOther)
	}
}
