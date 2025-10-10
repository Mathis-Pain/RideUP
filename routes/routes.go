package routes

import (
	"RideUP/handlers"
	"RideUP/utils"
	"net/http"
)

func InitRoutes() *http.ServeMux {

	mux := http.NewServeMux()

	// route Home
	// fonction anonyme utilisé pour faire dautre verification
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			utils.NotFoundHandler(w)
			return
		}
		handlers.HomeHandler(w, r)

	})
	// servir les fichiers static
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	return mux
}
