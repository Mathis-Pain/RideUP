package routes

import (
	"RideUP/handlers"
	"RideUP/middleware"
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
	mux.HandleFunc("/CreateAccount", handlers.RegisterHandler)
	mux.HandleFunc("/Connect", handlers.ConnectHandler)
	mux.HandleFunc("/Disconnect", handlers.DisconnectHandler)
	mux.Handle("/RideUp", middleware.AuthMiddleware(http.HandlerFunc(handlers.RideUpHandler)))
	mux.Handle("/NewRide", middleware.AuthMiddleware(http.HandlerFunc(handlers.NewRideHandler)))
	// servir les fichiers static
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	return mux
}
