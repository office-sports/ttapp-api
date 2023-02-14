package models

type Point struct {
	Id          int    `json:"id"`
	ScoreId     int    `json:"score_id"`
	IsHomePoint int    `json:"is_home_point"`
	IsAwayPoint int    `json:"is_away_point"`
	Time        string `json:"time"`
}

type PointPayload struct {
	GameId      int `json:"game_id"`
	IsHomePoint int `json:"is_home_point"`
	IsAwayPoint int `json:"is_away_point"`
}
