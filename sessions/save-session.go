package sessions

import (
	"database/sql"
	"encoding/json"

	"github.com/Mathis-Pain/RideUp/models"
)

// chemin d'acces a la db
const dbPath = "./data/RideUp.db"

// SaveSession sauvegarde ou met à jour une session
// Transforme les données Go en text Json pour pouvoir les sotcker dans la db
func SaveSession(db *sql.DB, session models.Session) error {
	dataJSON, err := json.Marshal(session.Data)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT OR REPLACE INTO sessions
		(id, user_id, data, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, session.ID, session.UserID, string(dataJSON), session.ExpiresAt, session.CreatedAt)

	return err
}

// ouvre juste la db pour sauvegarder la session via SaveSession
func SaveSessionToDB(session models.Session) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	return SaveSession(db, session)
}
