package main

import (
	builddb "RideUP/buildDB"
	"RideUP/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// initialisation de la bdd
	db, err := builddb.InitDB()
	if err != nil {
		fmt.Println("Erreur creation bdd :", err)
		return
	}
	defer db.Close()
	fmt.Println("Projet lancé, DB prête a l'emploi")

	//initialisation des routes
	mux := routes.InitRoutes()

	// Demarrage du serveur
	fmt.Println("serveur démarré sur http://localhost:5080...")
	if err := http.ListenAndServe(":5080", mux); err != nil {
		log.Fatal("Erreur serveur:", err)
	}

}
