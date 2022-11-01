package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/office-sports/ttapp-api/data"
	"github.com/office-sports/ttapp-api/models"
	"log"
	"net/http"
)

func GetOffices(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetOffices()
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
