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

func GetTournaments(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetTournaments()
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.Office))
		return
	} else if err != nil {
		log.Println("Unable to get metadata", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(md)
}

func GetLiveTournament(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetLiveTournament()
	if err == sql.ErrNoRows {
		json.NewEncoder(writer).Encode(new(models.Tournament))
		return
	} else if err != nil {
		log.Println("Unable to get tournament data", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(writer).Encode(md)
}

func GetTournamentSchedule(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	tid, _ := strconv.Atoi(params["id"])
	num, _ := strconv.Atoi(params["num"])
	if tid == 0 {
		log.Println("Invalid tournament id or limit")
		http.Error(writer, "Invalid tournament id or limit id", http.StatusBadRequest)
		return
	}
	schedule, err := data.GetTournamentSchedule(tid, num)
	if err != nil {
		log.Println("Unable to get tournament schedule", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(schedule)
}

func GetTournamentStandingsById(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := strconv.Atoi(params["id"])
	if id == 0 {
		log.Println("Invalid id")
		http.Error(writer, "Invalid tournament id", http.StatusBadRequest)
		return
	}
	standings, err := data.GetTournamentStandingsById(id)
	if err != nil {
		log.Println("Unable to get tournament standings", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(writer).Encode(standings)
}
