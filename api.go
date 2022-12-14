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
	// User
	//router.HandleFunc("/user/login", handlers.SignIn).Methods("POST", "OPTIONS")
	//router.HandleFunc("/user/changepass", handlers.ChangePass).Methods("POST", "OPTIONS")
	//router.HandleFunc("/user/{userId}/tournament/{tournamentId}/balance", handlers.GetUserTournamentBalance).Methods("GET")
	//// Bets
	//router.HandleFunc("/user/{userId}/tournament/{tournamentId}/bets", handlers.GetUserTournamentGamesBets).Methods("GET")
	//router.HandleFunc("/user/{userId}/tournament/{tournamentId}/tournamentbets", handlers.GetUserTournamentTournamentBets).Methods("GET")
	//router.HandleFunc("/user/{userId}/bets", handlers.GetUserGamesBets).Methods("GET")
	//router.HandleFunc("/bets/{gameId}", handlers.AddBet).Methods("POST", "OPTIONS")
	//router.HandleFunc("/tournamentbets/{outcomeId}", handlers.AddTournamentBet).Methods("POST", "OPTIONS")
	//router.HandleFunc("/bets/{gameId}", handlers.GetGameBets).Methods("GET")
	//
	//router.HandleFunc("/games/{gameId}/start", handlers.StartGame).Methods("POST", "OPTIONS")
	//router.HandleFunc("/games/{gameId}/finish", handlers.FinishGame).Methods("POST", "OPTIONS")

	//
	//router.HandleFunc("/tournament/{tournamentId}/gameranking", handlers.GetTournamentGameRanking).Methods("GET")
	//router.HandleFunc("/tournament/{tournamentId}/ranking", handlers.GetTournamentRanking).Methods("GET")
	//router.HandleFunc("/tournament/{tournamentId}/outcomes", handlers.GetTournamentOutcomes).Methods("GET")
	//router.HandleFunc("/tournament/{tournamentId}/start", handlers.StartTournament).Methods("POST", "OPTIONS")
	//router.HandleFunc("/tournament/meta", handlers.GetTournamentsMetadata).Methods("GET")
	//router.HandleFunc("/tournament/{tournamentId}/statistics", handlers.GetTournamentStatistics).Methods("GET")

	log.Printf("Listening on port %d", config.APIConfig.Port)
	log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.APIConfig.Port), router))
}
