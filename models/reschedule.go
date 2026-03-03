package models

type RescheduleRequest struct {
	ID                   int     `json:"id"`
	GameID               int     `json:"game_id"`
	RequesterID          int     `json:"requester_id"`
	RequesterName        string  `json:"requester_name"`
	ProposedDateFrom     string  `json:"proposed_date_from"`
	ProposedDateTo       string  `json:"proposed_date_to"`
	Status               string  `json:"status"`
	ParentID             *int    `json:"parent_id"`
	CreatedAt            string  `json:"created_at"`
	RespondedAt          *string `json:"responded_at"`
	HomePlayerID         int     `json:"home_player_id"`
	AwayPlayerID         int     `json:"away_player_id"`
	HomePlayerName       string  `json:"home_player_name"`
	AwayPlayerName       string  `json:"away_player_name"`
	CurrentDateOfMatch   string  `json:"current_date_of_match"`
	OriginalDateOfMatch  string  `json:"original_date_of_match"`
}
