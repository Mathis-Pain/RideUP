package handlers

import (
	"RideUP/utils"
	"html/template"
	"log"
	"net/http"
)

var RegistrationHtml = template.Must(template.ParseFiles("templates/registration.html"))

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := RegistrationHtml.Execute(w, nil)
	if err != nil {
		log.Printf("<RegistrationHtml.go> could not execute template <registration.html>: %v\n", err)
		utils.NotFoundHandler(w)
	}
}
