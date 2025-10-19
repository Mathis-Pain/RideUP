package handlers

import (
	"RideUP/models"
	"RideUP/sessions"
	"RideUP/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var NewEventHtml = template.Must(template.ParseFiles(
	"templates/newevent.html",
	"templates/inithtml/inithead.html",
	"templates/inithtml/initnav.html",
	"templates/inithtml/initfooter.html",
))

func GeocodeAddress(address string) (float64, float64, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Set("q", address)
	params.Set("format", "json")
	params.Set("limit", "1")

	// Création d'un client HTTP personnalisé
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Création de la requête
	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return 0, 0, err
	}

	// CRITIQUE : Ajout du User-Agent requis par Nominatim
	req.Header.Set("User-Agent", "RideUP/1.0 (ride-up-app)")

	// Exécution de la requête
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("erreur requête nominatim: %v", err)
	}
	defer resp.Body.Close()

	// Vérification du status code
	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("nominatim status: %d", resp.StatusCode)
	}

	var result models.NominatimResponseCoord
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, fmt.Errorf("erreur décodage JSON: %v", err)
	}

	if len(result) == 0 {
		return 0, 0, fmt.Errorf("adresse introuvable")
	}

	lat, err := strconv.ParseFloat(result[0].Lat, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("erreur parsing latitude: %v", err)
	}

	lon, err := strconv.ParseFloat(result[0].Lon, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("erreur parsing longitude: %v", err)
	}

	log.Printf("Géocodage réussi: %s -> [%f, %f]", address, lat, lon)
	return lat, lon, nil
}

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
		address := r.FormValue("address")  // Nouvelle saisie adresse
		latStr := r.FormValue("latitude")  // Champs hidden carte
		lonStr := r.FormValue("longitude") // Champs hidden carte

		if title == "" || dateStr == "" || timeStr == "" || (address == "" && (latStr == "" || lonStr == "")) {
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

		// Récupération des coordonnées
		var lat, lon float64
		if address != "" {
			// Géocodage adresse
			lat, lon, err = GeocodeAddress(address)
			if err != nil {
				log.Printf("Erreur géocodage : %v", err)
				http.Error(w, "Adresse introuvable", http.StatusBadRequest)
				return
			}
		} else {
			// Latitude/Longitude via clic sur carte
			lat, err = strconv.ParseFloat(latStr, 64)
			if err != nil {
				log.Printf("Erreur parsing latitude: %v", err)
				http.Redirect(w, r, "/NewEvent", http.StatusSeeOther)
				return
			}
			lon, err = strconv.ParseFloat(lonStr, 64)
			if err != nil {
				log.Printf("Erreur parsing longitude: %v", err)
				http.Redirect(w, r, "/NewEvent", http.StatusSeeOther)
				return
			}
		}

		// Insertion dans la base de données
		_, err = db.Exec(`
			INSERT INTO events (title, created_by, latitude, longitude, start_datetime)
			VALUES (?, ?, ?, ?, ?)`,
			title, session.UserID, lat, lon, startDatetime)
		if err != nil {
			log.Printf("ERREUR : insertion event: %v", err)
			http.Error(w, "Impossible de créer l'événement", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/RideUP", http.StatusSeeOther)
		return
	}

	// Affichage du formulaire
	err = NewEventHtml.Execute(w, data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template NewEventHtml: %v", err)
		utils.NotFoundHandler(w)
	}
}
