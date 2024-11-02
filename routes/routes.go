package routes

import (
	"github.com/gorilla/mux"
	"github.com/izsal/go-refresh-token/controllers"
	"github.com/izsal/go-refresh-token/middleware"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/refresh", controllers.RefreshToken).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)
	api.HandleFunc("/items", controllers.GetItems).Methods("GET")
	api.HandleFunc("/items/{id}", controllers.GetItem).Methods("GET")
	api.HandleFunc("/items", controllers.CreateItem).Methods("POST")
	api.HandleFunc("/items/{id}", controllers.UpdateItem).Methods("PUT")
	api.HandleFunc("/items/{id}", controllers.DeleteItem).Methods("DELETE")
}
