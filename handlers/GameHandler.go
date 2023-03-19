package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/office-sports/ttapp-api/data"
	"github.com/office-sports/ttapp-api/models"
	"log"
	"net/http"
	"strconv"
)

// GetGamesLive returns array of live games
func GetGamesLive(writer http.ResponseWriter, request *http.Request) {
	g, err := data.GetLiveGames()
	checkErrHTTP(err, writer, "Unable to get live games")

	json.NewEncoder(writer).Encode(g)
}

// GetTournamentGamesLive returns array of live games for requested tournament
func GetTournamentGamesLive(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	tid, _ := strconv.Atoi(params["id"])
	g, err := data.GetTournamentLiveGames(tid)
	checkErrHTTP(err, writer, "Unable to get live games for requested tournament")

	json.NewEncoder(writer).Encode(g)
}

// SaveGameScore handles manual score input (game finish)
func SaveGameScore(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	var gr models.GameSetResults
	err := json.NewDecoder(request.Body).Decode(&gr)

	checkErrHTTP(err, writer, "Unable to fetch game")

	data.SaveGameScore(gr)
}

// FinalizeGame ends set or set + game
func FinalizeGame(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	var sf models.SetFinal
	err := json.NewDecoder(request.Body).Decode(&sf)

	checkErrHTTP(err, writer, "Unable to fetch set score id")

	data.FinalizeGame(sf)
}

// ChangeServer flips the player id as server in the db
func ChangeServer(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	var p models.ChangeServerPayload
	err := json.NewDecoder(request.Body).Decode(&p)

	checkErrHTTP(err, writer, "Unable to fetch changing server payload")

	data.UpdateServer(p)

	serve, err := data.GetGameServe(p.GameId)
	checkErrHTTP(err, writer, "Unable to fetch server data")

	json.NewEncoder(writer).Encode(serve)
}

// GetGameModes returns array of game modes
func GetGameModes(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetGameModes()
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.GameMode))
		return
	}
	checkErrHTTP(err, writer, "Unable to fetch game modes")

	json.NewEncoder(writer).Encode(md)
}

// GetEloCache returns array of finished games with ELO values
func GetEloCache(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetEloCache()
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.EloCache))
		return
	}

	checkErrHTTP(err, writer, "Unable to fetch elo cache data")

	json.NewEncoder(writer).Encode(md)
}

// GetGameTimeline returns game timeline for requested id
func GetGameTimeline(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	tid, _ := strconv.Atoi(params["id"])
	timeline, err := data.GetGameTimeline(tid)
	checkErrHTTP(err, writer, "Unable to get game timeline")

	json.NewEncoder(writer).Encode(timeline)
}

// AnnounceGame sends announce message and updates db value for it
func AnnounceGame(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	gid, _ := strconv.Atoi(params["id"])
	data.AnnounceGame(gid)
}

// GetGameServe returns serve data for game
func GetGameServe(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	tid, _ := strconv.Atoi(params["id"])
	serve, err := data.GetGameServe(tid)
	checkErrHTTP(err, writer, "Unable to get game serve data")
	json.NewEncoder(writer).Encode(serve)
}

// GetGameById returns game for requested id
func GetGameById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid game id")
		http.Error(writer, "Invalid game id", http.StatusBadRequest)
		return
	}
	g, err := data.GetGameById(id)
	checkErrHTTP(err, writer, "Unable to get game by id")

	json.NewEncoder(writer).Encode(g)
}
