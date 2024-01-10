package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/office-sports/ttapp-api/data"
	"github.com/office-sports/ttapp-api/models"
	"log"
	"net/http"
	"strconv"
)

// GetTournamentPlayersStatistics returns array of tournaments
func GetTournamentPlayersStatistics(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid tournament id")
		http.Error(writer, "Invalid tournament id", http.StatusBadRequest)
		return
	}

	SetHeaders(writer)
	md, err := data.GetTournamentPlayersStatistics(id)
	if errors.Is(err, sql.ErrNoRows) {
		json.NewEncoder(writer).Encode(new(models.TournamentPlayerStatistics))
		return
	}
	checkErrHTTP(err, writer, "Unable to get tournament players statistics")

	json.NewEncoder(writer).Encode(md)
}

// GetTournamentsStatistics returns array of tournaments
func GetTournamentsStatistics(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetTournamentsStatistics()
	if errors.Is(err, sql.ErrNoRows) {
		json.NewEncoder(writer).Encode(new(models.TournamentStatistics))
		return
	}
	checkErrHTTP(err, writer, "Unable to get tournaments statistics")

	json.NewEncoder(writer).Encode(md)
}

// GetTournaments returns array of tournaments
func GetTournaments(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetTournaments()
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.Tournament))
		return
	}
	checkErrHTTP(err, writer, "Unable to get tournaments")

	json.NewEncoder(writer).Encode(md)
}

// GetLiveTournaments returns an array of all live tournaments
func GetLiveTournaments(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetLiveTournaments()
	if errors.Is(err, sql.ErrNoRows) {
		json.NewEncoder(writer).Encode(new(models.Tournament))
		return
	}
	checkErrHTTP(err, writer, "Unable to get live tournaments")

	json.NewEncoder(writer).Encode(md)
}

// GetTournamentById returns tournament data for requested id
func GetTournamentById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid tournament id")
		http.Error(writer, "Invalid tournament id", http.StatusBadRequest)
		return
	}
	tournament, err := data.GetTournamentById(id)
	checkErrHTTP(err, writer, "Unable to get tournament")

	json.NewEncoder(writer).Encode(tournament)
}

// GetTournamentSchedule returns tournament games ordered by date
func GetTournamentSchedule(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	tid, _ := strconv.Atoi(params["id"])
	num, _ := strconv.Atoi(params["num"])
	schedule, err := data.GetTournamentSchedule(tid, num)
	checkErrHTTP(err, writer, "Unable to get tournament schedule")

	json.NewEncoder(writer).Encode(schedule)
}

// GetTournamentGroupSchedule returns tournament games ordered by date
func GetTournamentGroupSchedule(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	tid, _ := strconv.Atoi(params["id"])
	gid, _ := strconv.Atoi(params["groupId"])
	schedule, err := data.GetTournamentGroupSchedule(tid, gid)
	checkErrHTTP(err, writer, "Unable to get tournament group schedule")

	json.NewEncoder(writer).Encode(schedule)
}

// GetTournamentResults returns array of finished games for requested tournament
func GetTournamentResults(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	tid, _ := strconv.Atoi(params["id"])
	num, _ := strconv.Atoi(params["num"])
	results, err := data.GetTournamentResults(tid, num)
	checkErrHTTP(err, writer, "Unable to get tournament results")

	json.NewEncoder(writer).Encode(results)
}

// GetTournamentGames returns array of all games for requested tournament
func GetTournamentGames(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	tid, _ := strconv.Atoi(params["id"])
	results, err := data.GetTournamentGames(tid)
	checkErrHTTP(err, writer, "Unable to get tournament games")

	json.NewEncoder(writer).Encode(results)
}

// GetTournamentStandingsById returns standings for requested tournament
func GetTournamentStandingsById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid tournament standings id")
		http.Error(writer, "Invalid tournament id", http.StatusBadRequest)
		return
	}
	standings, err := data.GetTournamentStandingsById(id)
	checkErrHTTP(err, writer, "Unable to get tournament standings")

	json.NewEncoder(writer).Encode(standings)
}

// GetTournamentPerformanceById returns players performance for requested tournament
func GetTournamentPerformanceById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid tournament performance id")
		http.Error(writer, "Invalid tournament id", http.StatusBadRequest)
		return
	}
	standings, err := data.GetTournamentPerformanceById(id)
	checkErrHTTP(err, writer, "Unable to get tournament performance")

	json.NewEncoder(writer).Encode(standings)
}

// GetTournamentInfo returns standings for requested tournament
func GetTournamentInfo(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid tournament info id")
		http.Error(writer, "Invalid tournament id", http.StatusBadRequest)
		return
	}
	info, err := data.GetTournamentInfo(id)
	checkErrHTTP(err, writer, "Unable to get tournament standings")

	json.NewEncoder(writer).Encode(info)
}

// GetTournamentLadders returns playoffs tournament ladders
func GetTournamentLadders(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid tournament ladder id")
		http.Error(writer, "Invalid tournament id", http.StatusBadRequest)
		return
	}
	ladders, err := data.GetTournamentLadders(id)
	checkErrHTTP(err, writer, "Unable to get tournament ladders")

	json.NewEncoder(writer).Encode(ladders)
}
