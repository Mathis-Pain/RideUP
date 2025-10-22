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
	"strconv"
	"time"
)

var NewEventHtml = template.Must(template.ParseFiles(
	"templates/newevent.html",
	"templates/inithtml/inithead.html",
	"templates/inithtml/initnav.html",
	"templates/inithtml/initfooter.html",
))

func NewEventHandler(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Printf("ERREUR : ParseForm: %v", err)
			http.Error(w, "Erreur lors de la lecture du formulaire", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		dateStr := r.FormValue("date")
		timeStr := r.FormValue("time")
		address := r.FormValue("address")  // Adresse (peut être vide si coords directes)
		latStr := r.FormValue("latitude")  // Champs hidden carte
		lonStr := r.FormValue("longitude") // Champs hidden carte

		if title == "" || dateStr == "" || timeStr == "" {
			http.Error(w, "Le titre, la date et l'heure sont obligatoires", http.StatusBadRequest)
			return
		}

		// Vérifier qu'on a soit une adresse, soit des coordonnées
		if address == "" && (latStr == "" || lonStr == "") {
			http.Error(w, "Veuillez saisir une adresse ou sélectionner un point sur la carte", http.StatusBadRequest)
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

		// Récupération des coordonnées et de l'adresse finale
		var lat, lon float64
		var finalAddress string

		if latStr != "" && lonStr != "" {
			// Coordonnées via carte (prioritaire)
			lat, err = strconv.ParseFloat(latStr, 64)
			if err != nil {
				log.Printf("Erreur parsing latitude: %v", err)
				http.Error(w, "Coordonnées invalides", http.StatusBadRequest)
				return
			}
			lon, err = strconv.ParseFloat(lonStr, 64)
			if err != nil {
				log.Printf("Erreur parsing longitude: %v", err)
				http.Error(w, "Coordonnées invalides", http.StatusBadRequest)
				return
			}
			// Si une adresse est fournie (via géocodage inversé JS), on l'utilise
			if address != "" {
				finalAddress = address
			} else {
				finalAddress = fmt.Sprintf("Lat: %.6f, Lon: %.6f", lat, lon)
			}

		} else if address != "" {
			// Géocodage adresse uniquement si pas de coordonnées
			lat, lon, err = utils.GeocodeAddress(address)
			if err != nil {
				log.Printf("Erreur géocodage : %v", err)
				http.Error(w, "Adresse introuvable. Veuillez vérifier l'orthographe ou utiliser la carte.", http.StatusBadRequest)
				return
			}
			finalAddress = address
		}

		log.Printf("INFO - Création événement: %s à %s [%.6f, %.6f]", title, finalAddress, lat, lon)

		// Insertion dans la base de données avec l'adresse
		_, err = db.Exec(`
			INSERT INTO events (title, created_by, latitude, longitude, address, start_datetime)
			VALUES (?, ?, ?, ?, ?, ?)`,
			title, session.UserID, lat, lon, finalAddress, startDatetime)
		if err != nil {
			log.Printf("ERREUR : insertion event: %v", err)
			http.Error(w, "Impossible de créer l'événement", http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Événement créé avec succès: %s", title)
		http.Redirect(w, r, "/RideUp", http.StatusSeeOther)
		return
	}

	// Affichage du formulaire
	err = NewEventHtml.Execute(w, data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template NewEventHtml: %v", err)
		utils.NotFoundHandler(w)
	}
}
