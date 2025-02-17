package routing

import (
	"boilerplate/controllers"

	"github.com/gorilla/mux"
)

func User(router *mux.Router) {
	router.HandleFunc("", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	authenticated := router.PathPrefix("").Subrouter()
	authenticated.Use(AuthMiddleware)
	authenticated.HandleFunc("/me", controllers.GetMe).Methods("GET")
	authenticated.HandleFunc("/{id}", controllers.GetUser).Methods("GET")
}