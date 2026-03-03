package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/office-sports/ttapp-api/data"
)

type RescheduleRequestBody struct {
	ProposedDateFrom string `json:"proposed_date_from"`
	ProposedDateTo   string `json:"proposed_date_to"`
}

// CreateRescheduleRequest handles POST /games/{id}/reschedule
func CreateRescheduleRequest(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	claims := r.Context().Value("claims").(*AuthClaims)
	if claims == nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	gameID, err := strconv.Atoi(params["id"])
	if err != nil || gameID == 0 {
		http.Error(w, `{"error":"Invalid game id"}`, http.StatusBadRequest)
		return
	}

	var body RescheduleRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ProposedDateFrom == "" {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}
	if body.ProposedDateTo == "" {
		body.ProposedDateTo = body.ProposedDateFrom
	}

	// Reject past dates
	today := time.Now().Format("2006-01-02")
	if body.ProposedDateFrom < today {
		http.Error(w, `{"error":"Proposed date cannot be in the past"}`, http.StatusBadRequest)
		return
	}

	// Reject if a pending request already exists for this game
	hasPending, err := data.HasPendingRescheduleForGame(gameID)
	if err != nil {
		log.Printf("[RESCHEDULE] HasPendingRescheduleForGame gameID=%d: %v", gameID, err)
		http.Error(w, `{"error":"Failed to check existing requests"}`, http.StatusInternalServerError)
		return
	}
	if hasPending {
		http.Error(w, `{"error":"A reschedule request is already pending for this match"}`, http.StatusConflict)
		return
	}

	id, err := data.CreateRescheduleRequest(gameID, claims.PlayerID, body.ProposedDateFrom, body.ProposedDateTo, nil)
	if err != nil {
		log.Printf("[RESCHEDULE] CreateRescheduleRequest gameID=%d playerID=%d: %v", gameID, claims.PlayerID, err)
		http.Error(w, `{"error":"Failed to create reschedule request"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// GetRescheduleRequestsForPlayer handles GET /players/{id}/reschedule
func GetRescheduleRequestsForPlayer(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	claims := r.Context().Value("claims").(*AuthClaims)
	if claims == nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	playerID, err := strconv.Atoi(params["id"])
	if err != nil || playerID == 0 {
		http.Error(w, `{"error":"Invalid player id"}`, http.StatusBadRequest)
		return
	}

	// Only allow fetching own requests
	if claims.PlayerID != playerID {
		http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
		return
	}

	requests, err := data.GetRescheduleRequestsForPlayer(playerID)
	if err != nil {
		log.Printf("[RESCHEDULE] GetRescheduleRequestsForPlayer playerID=%d: %v", playerID, err)
		http.Error(w, `{"error":"Failed to get reschedule requests"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requests)
}

// AcceptRescheduleRequest handles POST /reschedule/{id}/accept
func AcceptRescheduleRequest(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	claims := r.Context().Value("claims").(*AuthClaims)
	if claims == nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	reqID, err := strconv.Atoi(params["id"])
	if err != nil || reqID == 0 {
		http.Error(w, `{"error":"Invalid reschedule request id"}`, http.StatusBadRequest)
		return
	}

	req, err := data.GetRescheduleRequestByID(reqID)
	if err != nil {
		log.Printf("[RESCHEDULE] AcceptRescheduleRequest GetRescheduleRequestByID id=%d: %v", reqID, err)
		http.Error(w, `{"error":"Reschedule request not found"}`, http.StatusNotFound)
		return
	}

	// Only the opponent (non-requester participant) can accept
	if claims.PlayerID == req.RequesterID {
		http.Error(w, `{"error":"Cannot accept your own request"}`, http.StatusForbidden)
		return
	}
	if claims.PlayerID != req.HomePlayerID && claims.PlayerID != req.AwayPlayerID {
		http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
		return
	}

	// For range proposals, accept body may contain a confirmed_date
	var body struct {
		ConfirmedDate string `json:"confirmed_date"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	agreedDate := req.ProposedDateFrom
	if body.ConfirmedDate != "" {
		agreedDate = body.ConfirmedDate
	}

	if err := data.UpdateRescheduleStatus(reqID, "accepted"); err != nil {
		log.Printf("[RESCHEDULE] AcceptRescheduleRequest UpdateRescheduleStatus id=%d: %v", reqID, err)
		http.Error(w, `{"error":"Failed to accept"}`, http.StatusInternalServerError)
		return
	}

	// Update the game date with the agreed date
	if err := data.UpdateGameDateOfMatch(req.GameID, agreedDate); err != nil {
		log.Printf("[RESCHEDULE] AcceptRescheduleRequest UpdateGameDateOfMatch gameID=%d: %v", req.GameID, err)
		http.Error(w, `{"error":"Failed to update game date"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}

// DeclineRescheduleRequest handles POST /reschedule/{id}/decline
func DeclineRescheduleRequest(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	claims := r.Context().Value("claims").(*AuthClaims)
	if claims == nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	reqID, err := strconv.Atoi(params["id"])
	if err != nil || reqID == 0 {
		http.Error(w, `{"error":"Invalid reschedule request id"}`, http.StatusBadRequest)
		return
	}

	req, err := data.GetRescheduleRequestByID(reqID)
	if err != nil {
		log.Printf("[RESCHEDULE] DeclineRescheduleRequest GetRescheduleRequestByID id=%d: %v", reqID, err)
		http.Error(w, `{"error":"Reschedule request not found"}`, http.StatusNotFound)
		return
	}

	if claims.PlayerID != req.HomePlayerID && claims.PlayerID != req.AwayPlayerID {
		http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
		return
	}

	if err := data.UpdateRescheduleStatus(reqID, "declined"); err != nil {
		log.Printf("[RESCHEDULE] DeclineRescheduleRequest UpdateRescheduleStatus id=%d: %v", reqID, err)
		http.Error(w, `{"error":"Failed to decline"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "declined"})
}

// RevokeRescheduleRequest handles POST /reschedule/{id}/revoke
// Restores the original game date and marks the request as revoked.
func RevokeRescheduleRequest(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	claims := r.Context().Value("claims").(*AuthClaims)
	if claims == nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	reqID, err := strconv.Atoi(params["id"])
	if err != nil || reqID == 0 {
		http.Error(w, `{"error":"Invalid reschedule request id"}`, http.StatusBadRequest)
		return
	}

	req, err := data.GetRescheduleRequestByID(reqID)
	if err != nil {
		log.Printf("[RESCHEDULE] RevokeRescheduleRequest GetRescheduleRequestByID id=%d: %v", reqID, err)
		http.Error(w, `{"error":"Reschedule request not found"}`, http.StatusNotFound)
		return
	}

	if req.Status != "accepted" {
		http.Error(w, `{"error":"Only accepted requests can be revoked"}`, http.StatusBadRequest)
		return
	}

	if claims.PlayerID != req.HomePlayerID && claims.PlayerID != req.AwayPlayerID {
		http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
		return
	}

	if err := data.UpdateRescheduleStatus(reqID, "revoked"); err != nil {
		log.Printf("[RESCHEDULE] RevokeRescheduleRequest UpdateRescheduleStatus id=%d: %v", reqID, err)
		http.Error(w, `{"error":"Failed to revoke"}`, http.StatusInternalServerError)
		return
	}

	// Restore original game date (may be empty if it was unset)
	if err := data.UpdateGameDateOfMatch(req.GameID, req.OriginalDateOfMatch); err != nil {
		log.Printf("[RESCHEDULE] RevokeRescheduleRequest UpdateGameDateOfMatch gameID=%d: %v", req.GameID, err)
		http.Error(w, `{"error":"Failed to restore original date"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "revoked"})
}

func CounterRescheduleRequest(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)
	claims := r.Context().Value("claims").(*AuthClaims)
	if claims == nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	reqID, err := strconv.Atoi(params["id"])
	if err != nil || reqID == 0 {
		http.Error(w, `{"error":"Invalid reschedule request id"}`, http.StatusBadRequest)
		return
	}

	req, err := data.GetRescheduleRequestByID(reqID)
	if err != nil {
		log.Printf("[RESCHEDULE] CounterRescheduleRequest GetRescheduleRequestByID id=%d: %v", reqID, err)
		http.Error(w, `{"error":"Reschedule request not found"}`, http.StatusNotFound)
		return
	}

	if claims.PlayerID != req.HomePlayerID && claims.PlayerID != req.AwayPlayerID {
		http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
		return
	}

	var body RescheduleRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ProposedDateFrom == "" {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}
	if body.ProposedDateTo == "" {
		body.ProposedDateTo = body.ProposedDateFrom
	}

	// Overwrite the single row for this game with the counter-proposal
	if err := data.UpdateRescheduleByGameID(req.GameID, claims.PlayerID, body.ProposedDateFrom, body.ProposedDateTo); err != nil {
		log.Printf("[RESCHEDULE] CounterRescheduleRequest UpdateRescheduleByGameID gameID=%d: %v", req.GameID, err)
		http.Error(w, `{"error":"Failed to save counter proposal"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": req.GameID})
}
