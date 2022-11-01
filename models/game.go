package models

type GameResult struct {
	MatchId        int        `json:"match_id"`
	GroupName      string     `json:"group_name"`
	DateOfMatch    string     `json:"date_of_match"`
	DatePlayed     *string    `json:"date_played"`
	HomePlayerId   int        `json:"home_player_id"`
	AwayPlayerId   int        `json:"away_player_id"`
	HomePlayerName string     `json:"home_player_name"`
	AwayPlayerName string     `json:"away_player_name"`
	WinnerId       int        `json:"winner_id"`
	HomeScoreTotal int        `json:"home_score_total"`
	AwayScoreTotal int        `json:"away_score_total"`
	HomeElo        int        `json:"home_elo"`
	AwayElo        int        `json:"away_elo"`
	HomeEloDiff    int        `json:"home_elo_diff"`
	AwayEloDiff    int        `json:"away_elo_diff"`
	SetScores      []SetScore `json:"scores"`
}

type SetScore struct {
	Set  int `json:"set"`
	Home int `json:"home"`
	Away int `json:"away"`
}

type GameResultSetScores struct {
	S1hp *int `json:"s1hp"`
	S1ap *int `json:"s1ap"`
	S2hp *int `json:"s2hp"`
	S2ap *int `json:"s2ap"`
	S3hp *int `json:"s3hp"`
	S3ap *int `json:"s3ap"`
	S4hp *int `json:"s4hp"`
	S4ap *int `json:"s4ap"`
	S5hp *int `json:"s5hp"`
	S5ap *int `json:"s5ap"`
}
