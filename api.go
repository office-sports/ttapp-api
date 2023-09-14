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

	// office routes
	router.HandleFunc("/offices", handlers.GetOffices).Methods("GET")

	// player routes
	router.HandleFunc("/players", handlers.GetPlayers).Methods("GET")
	router.HandleFunc("/players/{id}", handlers.GetPlayerById).Methods("GET")
	router.HandleFunc("/players/{id}/results", handlers.GetPlayerResultsById).Methods("GET")
	router.HandleFunc("/players/{id}/schedule", handlers.GetPlayerScheduleById).Methods("GET")

	// tournament routes
	router.HandleFunc("/tournaments", handlers.GetTournaments).Methods("GET")
	router.HandleFunc("/tournaments/{id}/statistics", handlers.GetTournamentsStatistics).Methods("GET")
	router.HandleFunc("/tournaments/{id}/players_statistics", handlers.GetTournamentPlayersStatistics).Methods("GET")
	router.HandleFunc("/tournaments/live", handlers.GetLiveTournaments).Methods("GET")
	router.HandleFunc("/tournaments/{id}", handlers.GetTournamentById).Methods("GET")
	router.HandleFunc("/tournaments/{id}/schedule/{num}", handlers.GetTournamentSchedule).Methods("GET")
	router.HandleFunc("/tournaments/{id}/results/{num}", handlers.GetTournamentResults).Methods("GET")
	router.HandleFunc("/tournaments/{id}/games", handlers.GetTournamentGames).Methods("GET")
	router.HandleFunc("/tournaments/{id}/standings", handlers.GetTournamentStandingsById).Methods("GET")
	router.HandleFunc("/tournaments/{id}/performance", handlers.GetTournamentPerformanceById).Methods("GET")
	router.HandleFunc("/tournaments/{id}/ladders", handlers.GetTournamentLadders).Methods("GET")
	router.HandleFunc("/tournaments/{id}/live_games", handlers.GetTournamentGamesLive).Methods("GET")
	router.HandleFunc("/tournaments/{id}/info", handlers.GetTournamentInfo).Methods("GET")

	// game routes
	router.HandleFunc("/games/live", handlers.GetGamesLive).Methods("GET")
	router.HandleFunc("/games/save", handlers.SaveGameScore).Methods("POST")
	router.HandleFunc("/games/finalize", handlers.FinalizeGame).Methods("POST")
	router.HandleFunc("/games/changeserver", handlers.ChangeServer).Methods("POST")
	router.HandleFunc("/games/modes", handlers.GetGameModes).Methods("GET")
	router.HandleFunc("/games/{id}/details", handlers.GetGameTimeline).Methods("GET")
	router.HandleFunc("/games/{id}/announce", handlers.AnnounceGame).Methods("GET")
	router.HandleFunc("/games/{id}/serve", handlers.GetGameServe).Methods("GET")
	router.HandleFunc("/games/{id}/elo", handlers.UpdateGameElo).Methods("GET")
	router.HandleFunc("/games/{id}", handlers.GetGameById).Methods("GET")

	// scoring routes
	router.HandleFunc("/points/add", handlers.AddPoint).Methods("POST")
	router.HandleFunc("/points/del", handlers.DelPoint).Methods("POST")

	// additional routes
	router.HandleFunc("/leaders", handlers.GetLeaders).Methods("GET")

	log.Printf("Listening on port %d", config.APIConfig.Port)
	log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.APIConfig.Port), router))
}
