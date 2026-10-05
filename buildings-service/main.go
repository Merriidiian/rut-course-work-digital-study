package main

import (
	"net/http"

	"university/internal/buildings"
	"university/internal/server"
)

func main() {
	database := server.OpenDatabase()
	defer database.Close()

	handler := buildings.NewHandler(database)
	router := http.NewServeMux()
	handler.Register(router)
	server.Serve(router)
}
