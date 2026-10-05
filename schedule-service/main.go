package main

import (
	"net/http"

	"university/schedule/internal/schedule"
	"university/schedule/internal/server"
)

func main() {
	database := server.OpenDatabase()
	defer database.Close()

	handler := schedule.NewHandler(database)
	router := http.NewServeMux()
	handler.Register(router)
	server.Serve(router)
}
