package main

import (
	"net/http"

	"university/rooms/internal/rooms"
	"university/rooms/internal/server"
)

func main() {
	database := server.OpenDatabase()
	defer database.Close()

	handler := rooms.NewHandler(database)
	router := http.NewServeMux()
	handler.Register(router)
	server.Serve(router)
}
