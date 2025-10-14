package handlers

import (
	"RideUP/utils"
	"html/template"
	"log"
	"net/http"
)

var ConnectionHtml = template.Must(template.ParseFiles("templates/connection.html"))

func ConnectHandler(w http.ResponseWriter, r *http.Request) {
	err := ConnectionHtml.Execute(w, nil)
	if err != nil {
		log.Printf("<Connect.go> could not execute template <connection.html>: %v\n", err)
		utils.NotFoundHandler(w)
	}
}
