package builddb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func InitDB() (*sql.DB, error) {
	// 1. Ouvrir ou créer la base de données
	db, err := sql.Open("sqlite3", "/data/RideUp.db")
	if err != nil {
		log.Fatal("Impossible d'ouvrir la DB:", err)
	}
	defer db.Close()

	// 2. Activer les foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Impossible d'activer les foreign keys:", err)
	}

	// 3. Lire le fichier SQL
	schema, err := os.ReadFile("/data/schemaRideUp.sql")
	if err != nil {
		log.Fatal("Impossible de lire le fichier SQL:", err)
	}

	// 4. Exécuter le script SQL
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal("Erreur lors de la création des tables:", err)
	}

	fmt.Println("Base de données créée avec succès !")
	return db, nil
}
