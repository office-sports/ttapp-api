package handlers

import (
	"encoding/json"
	"github.com/office-sports/ttapp-api/data"
	"github.com/office-sports/ttapp-api/models"
	"log"
	"net/http"
)

// AddPoint handles adding point to the database table
func AddPoint(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	var p models.PointPayload
	err := json.NewDecoder(request.Body).Decode(&p)

	s, err := data.GetScoresByGameId(p.GameId)
	checkErrHTTP(err, writer, "Unable to get scores for game")

	// check if we have sets data for this game
	if len(s) == 0 {
		var z int = 0

		// no sets data, insert new one
		lid, err := data.InsertSetScore(p.GameId, 1, &z, &z)
		if err != nil {
			log.Println("Error inserting new score", err)
		}

		// update current set in game
		data.UpdateGameCurrentSet(p.GameId, 1)

		// add point with reference to the score id
		_, err = data.InsertPoint(int(lid), p.IsHomePoint, p.IsAwayPoint)
		if err != nil {
			log.Println("Error inserting new point", err)
		}
	} else {
		// get current set for the game
		g, _ := data.GetGameById(p.GameId)

		// add point with reference to the score id
		_, err = data.InsertPoint(*g.CurrentSetId, p.IsHomePoint, p.IsAwayPoint)
		if err != nil {
			log.Println("Error inserting new point", err)
		}
	}

	serve, err := data.GetGameServe(p.GameId)
	checkErrHTTP(err, writer, "Unable to get game serve data")

	json.NewEncoder(writer).Encode(serve)
}

// DelPoint handles removing point from the database table
func DelPoint(writer http.ResponseWriter, request *http.Request) {
	SetHeaders(writer)
	var p models.PointPayload
	err := json.NewDecoder(request.Body).Decode(&p)

	maxPid, err := data.GetMaxPoint(p.GameId, p.IsHomePoint, p.IsAwayPoint)

	checkErrHTTP(err, writer, "Unable to get max point for game")

	if maxPid != 0 {
		data.DeletePointById(maxPid)
	}

	serve, err := data.GetGameServe(p.GameId)
	checkErrHTTP(err, writer, "Unable to get game serve data")

	json.NewEncoder(writer).Encode(serve)
}
