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

func GetGameModes(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetGameModes()
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.GameMode))
		return
	} else if err != nil {
		log.Println("Unable to get game modes", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(md)
}

func GetEloCache(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetEloCache()
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.EloCache))
		return
	} else if err != nil {
		log.Println("Unable to get elo cache", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(md)
}

func GetGameTimeline(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	tid, _ := strconv.Atoi(params["id"])
	timeline, err := data.GetGameTimeline(tid)
	if err != nil {
		log.Println("Unable to get game timeline", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(timeline)
}

func GetGameServe(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	tid, _ := strconv.Atoi(params["id"])
	serve, err := data.GetGameServe(tid)
	if err != nil {
		log.Println("Unable to get game serve data", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(serve)
}

func GetGameById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid game id")
		http.Error(writer, "Invalid game id", http.StatusBadRequest)
		return
	}
	g, err := data.GetGameById(id)
	if err != nil {
		log.Println("Unable to get game", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(g)
}

func SaveGame(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	var gr models.GameSetResults
	err := json.NewDecoder(request.Body).Decode(&gr)

	if err != nil {
		log.Println("Unable to get game id", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	data.Save(gr)
}

func ChangeServer(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	var p models.ChangeServerPayload
	err := json.NewDecoder(request.Body).Decode(&p)

	if err != nil {
		log.Println("Unable to get game id", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	data.UpdateServer(p)

	serve, err := data.GetGameServe(p.GameId)
	if err != nil {
		log.Println("Unable to get game serve data", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(serve)
}
