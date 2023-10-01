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

// GetPlayers return all players
func GetPlayers(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetPlayers()
	if errors.Is(err, sql.ErrNoRows) {
		err := json.NewEncoder(writer).Encode(new(models.Player))
		checkErrSimple(err)
		return
	}
	checkErrHTTP(err, writer, "Unable to get players data")

	json.NewEncoder(writer).Encode(md)
}

// GetPlayerById returns single player by given id
func GetPlayerById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid player id")
		http.Error(writer, "Invalid player id", http.StatusBadRequest)
		return
	}
	player, err := data.GetPlayerById(id)
	checkErrHTTP(err, writer, "Unable to get player")

	json.NewEncoder(writer).Encode(player)
}

// GetPlayerResultsById returns array of player results for given player id
func GetPlayerResultsById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid player id")
		http.Error(writer, "Invalid player id", http.StatusBadRequest)
		return
	}
	player, err := data.GetPlayerGamesById(id, 1)
	checkErrHTTP(err, writer, "Unable to get player results")

	json.NewEncoder(writer).Encode(player)
}

// GetPlayerScheduleById returns array of player scheduled games for given player id
func GetPlayerScheduleById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid player id")
		http.Error(writer, "Invalid player id", http.StatusBadRequest)
		return
	}
	player, err := data.GetPlayerGamesById(id, 0)
	checkErrHTTP(err, writer, "Unable to get player schedule")

	json.NewEncoder(writer).Encode(player)
}

// GetPlayerOpponentsById returns array of player opponents
func GetPlayerOpponentsById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid player id")
		http.Error(writer, "Invalid player id", http.StatusBadRequest)
		return
	}
	player, err := data.GetPlayerOpponentsById(id)
	checkErrHTTP(err, writer, "Unable to get player opponents")

	json.NewEncoder(writer).Encode(player)
}

// GetLeaders returns all-time leaders
func GetLeaders(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetLeaders()
	if errors.Is(err, sql.ErrNoRows) {
		json.NewEncoder(writer).Encode(new(models.Leader))
		return
	}
	checkErrHTTP(err, writer, "Unable to get leaders data")

	json.NewEncoder(writer).Encode(md)
}
