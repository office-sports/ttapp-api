package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/office-sports/ttapp-api/data"
	"github.com/office-sports/ttapp-api/models"
	"net/http"
)

func GetOffices(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	md, err := data.GetOffices()
	if err == sql.ErrNoRows {
		err := json.NewEncoder(writer).Encode(new(models.Office))
		checkErrSimple(err)
		return
	}

	checkErrHTTP(err, writer, "Unable to get metadata")

	err = json.NewEncoder(writer).Encode(md)
	checkErrSimple(err)
}
