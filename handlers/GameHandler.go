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
