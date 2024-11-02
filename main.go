package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/izsal/go-refresh-token/config"
	"github.com/izsal/go-refresh-token/routes"
)

func main() {
	config.LoadConfig()
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	log.Println("Server Running on port 8080")
	http.ListenAndServe(":8080", r)
}
