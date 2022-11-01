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

func GetPlayers(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetPlayers()
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.Office))
		return
	} else if err != nil {
		log.Println("Unable to get players data", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(md)
}

func GetPlayerById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid id")
		http.Error(writer, "Invalid player id", http.StatusBadRequest)
		return
	}
	player, err := data.GetPlayerById(id)
	if err != nil {
		log.Println("Unable to get player", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(player)
}

func GetPlayerResultsById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid id")
		http.Error(writer, "Invalid player id", http.StatusBadRequest)
		return
	}
	player, err := data.GetPlayerGamesById(id, 1)
	if err != nil {
		log.Println("Unable to get player results", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(player)
}

func GetPlayerScheduleById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid id")
		http.Error(writer, "Invalid player id", http.StatusBadRequest)
		return
	}
	player, err := data.GetPlayerGamesById(id, 0)
	if err != nil {
		log.Println("Unable to get player results", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(player)
}
