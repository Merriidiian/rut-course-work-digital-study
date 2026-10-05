package main

import (
	"net/http"

	"university/buildings/internal/buildings"
	"university/buildings/internal/server"
)

func main() {
	database := server.OpenDatabase()
	defer database.Close()

	handler := buildings.NewHandler(database)
	router := http.NewServeMux()
	handler.Register(router)
	server.Serve(router)
}
