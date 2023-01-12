package main

import (
	"fmt"
	"github.com/office-sports/ttapp-api/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/office-sports/ttapp-api/models"
)

func main() {
	var configFileParam string

	if len(os.Args) > 1 {
		configFileParam = os.Args[1]
	}

	config, err := models.GetConfig(configFileParam)
	if err != nil {
		panic(err)
	}
	models.InitDB(config.GetMysqlConnectionString())

	router := mux.NewRouter()
	router.HandleFunc("/offices", handlers.GetOffices).Methods("GET")

	router.HandleFunc("/players", handlers.GetPlayers).Methods("GET")
	router.HandleFunc("/players/{id}", handlers.GetPlayerById).Methods("GET")
	router.HandleFunc("/players/{id}/results", handlers.GetPlayerResultsById).Methods("GET")
	router.HandleFunc("/players/{id}/schedule", handlers.GetPlayerScheduleById).Methods("GET")

	router.HandleFunc("/tournaments", handlers.GetTournaments).Methods("GET")
	router.HandleFunc("/tournaments/live", handlers.GetLiveTournament).Methods("GET")
	router.HandleFunc("/tournaments/{id}/schedule/{num}", handlers.GetTournamentSchedule).Methods("GET")
	router.HandleFunc("/tournaments/{id}/standings", handlers.GetTournamentStandingsById).Methods("GET")

	router.HandleFunc("/games/modes", handlers.GetGameModes).Methods("GET")

	log.Printf("Listening on port %d", config.APIConfig.Port)
	log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.APIConfig.Port), router))
}
