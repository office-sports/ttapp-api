package models

type Score struct {
	Id        int `json:"id"`
	GameId    int `json:"game_id"`
	SetNumber int `json:"set_number"`
	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`
}
