package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/office-sports/ttapp-api/data"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

type LoginRequest struct {
	SlackName string `json:"slack_name"`
	Password  string `json:"password"`
}

type LoginResponse struct {
	Token                string `json:"token"`
	PlayerID             int    `json:"player_id"`
	Name                 string `json:"name"`
	Nickname             string `json:"nickname"`
	ProfilePic           string `json:"profile_pic_url"`
	MustChangePassword   bool   `json:"must_change_password"`
}

type ProfileUpdateRequest struct {
	Nickname  string `json:"nickname"`
	Password  string `json:"password,omitempty"`
}

type AuthClaims struct {
	PlayerID  int    `json:"player_id"`
	SlackName string `json:"slack_name"`
	Name      string `json:"name"`
	jwt.RegisteredClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[AUTH] Failed to decode login request: %v", err)
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	log.Printf("[AUTH] Login attempt for username: %s", req.SlackName)

	// Add @ symbol to slack name (database stores with @)
	slackName := "@" + strings.TrimPrefix(req.SlackName, "@")
	log.Printf("[AUTH] Looking up player with slack_name: %s", slackName)

	player, err := data.GetPlayerBySlackName(slackName)
	if err != nil {
		log.Printf("[AUTH] Player not found with slack_name %s: %v", slackName, err)
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	log.Printf("[AUTH] Player found: ID=%d, Name=%s", player.ID, player.Name)

	storedHash, mustChange, err := data.GetPlayerPasswordHashWithFlag(player.ID)
	if err != nil || storedHash == "" {
		log.Printf("[AUTH] Password hash not found for player ID %d: %v", player.ID, err)
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	log.Printf("[AUTH] Comparing passwords for player ID: %d", player.ID)
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
		log.Printf("[AUTH] Password mismatch for player ID %d", player.ID)
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	log.Printf("[AUTH] Password verified for player ID: %d", player.ID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{
		PlayerID:  player.ID,
		SlackName: player.SlackName,
		Name:      player.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("[AUTH] Failed to generate JWT token for player ID %d: %v", player.ID, err)
		http.Error(w, `{"error":"Token generation failed"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("[AUTH] JWT token generated successfully for player ID: %d", player.ID)

	response := LoginResponse{
		Token:              tokenString,
		PlayerID:           player.ID,
		Name:               player.Name,
		Nickname:           getNicknameOrEmpty(player.Nickname),
		ProfilePic:         player.ProfilePicUrl,
		MustChangePassword: mustChange,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Printf("[AUTH] Login successful for player ID: %d (%s)", player.ID, player.Name)
}

// Helper function to handle nil nickname
func getNicknameOrEmpty(nickname *string) string {
	if nickname == nil {
		return ""
	}
	return *nickname
}

func SetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := r.Context().Value("claims").(*AuthClaims)
	if claims == nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, `{"error":"Password cannot be empty"}`, http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"Failed to process password"}`, http.StatusInternalServerError)
		return
	}

	if err := data.SetPlayerPasswordHash(claims.PlayerID, string(hash)); err != nil {
		http.Error(w, `{"error":"Failed to update password"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated"})
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := r.Context().Value("claims").(*AuthClaims)
	if claims == nil {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req ProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	if req.Nickname != "" {
		if err := data.UpdatePlayerNickname(claims.PlayerID, req.Nickname); err != nil {
			http.Error(w, `{"error":"Failed to update nickname"}`, http.StatusInternalServerError)
			return
		}
	}

	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, `{"error":"Failed to process password"}`, http.StatusInternalServerError)
			return
		}
		if err := data.SetPlayerPasswordHash(claims.PlayerID, string(hash)); err != nil {
			http.Error(w, `{"error":"Failed to update password"}`, http.StatusInternalServerError)
			return
		}
		// Clear the must_change_password flag after password is changed
		if err := data.ClearMustChangePasswordFlag(claims.PlayerID); err != nil {
			http.Error(w, `{"error":"Failed to update password status"}`, http.StatusInternalServerError)
			return
		}
	}

	player, err := data.GetPlayerById(claims.PlayerID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch updated profile"}`, http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token:     "",
		PlayerID:  player.ID,
		Name:      player.Name,
		Nickname:  getNicknameOrEmpty(player.Nickname),
		ProfilePic: player.ProfilePicUrl,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := r.Context().Value("claims").(*AuthClaims)
	if claims == nil {
		http.Error(w, `{"error":"Invalid token"}`, http.StatusUnauthorized)
		return
	}

	player, err := data.GetPlayerById(claims.PlayerID)
	if err != nil {
		http.Error(w, `{"error":"Player not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(player)
}
