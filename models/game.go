package models

type Game struct {
	TournamentId   int    `json:"tournament_id"`
	OfficeId       int    `json:"office_id"`
	MatchId        int    `json:"match_id"`
	GroupName      string `json:"group_name"`
	DateOfMatch    string `json:"date_of_match"`
	HomePlayerId   int    `json:"home_player_id"`
	AwayPlayerId   int    `json:"away_player_id"`
	HomePlayerName string `json:"home_player_name"`
	AwayPlayerName string `json:"away_player_name"`
	HomeScoreTotal int    `json:"home_score_total"`
	AwayScoreTotal int    `json:"away_score_total"`
}

type GameResult struct {
	MatchId           int        `json:"match_id"`
	MaxSets           int        `json:"max_sets"`
	WinsRequired      int        `json:"wins_required"`
	GroupName         string     `json:"group_name"`
	TournamentId      int        `json:"tournament_id"`
	OfficeId          int        `json:"office_id"`
	DateOfMatch       string     `json:"date_of_match"`
	DatePlayed        *string    `json:"date_played"`
	HomePlayerId      int        `json:"home_player_id"`
	AwayPlayerId      int        `json:"away_player_id"`
	HomePlayerName    string     `json:"home_player_name"`
	AwayPlayerName    string     `json:"away_player_name"`
	WinnerId          int        `json:"winner_id"`
	HomeScoreTotal    int        `json:"home_score_total"`
	AwayScoreTotal    int        `json:"away_score_total"`
	IsWalkover        int        `json:"is_walkover"`
	IsFinished        int        `json:"is_finished"`
	HomeElo           *int       `json:"home_elo"`
	AwayElo           *int       `json:"away_elo"`
	NewHomeElo        *int       `json:"new_home_elo"`
	NewAwayElo        *int       `json:"new_away_elo"`
	HomeEloDiff       *int       `json:"home_elo_diff"`
	AwayEloDiff       *int       `json:"away_elo_diff"`
	HasPoints         *int       `json:"has_points"`
	CurrentHomePoints *int       `json:"current_home_points"`
	CurrentAwayPoints *int       `json:"current_away_points"`
	CurrentSet        int        `json:"current_set"`
	CurrentSetId      *int       `json:"current_set_id"`
	Announced         int        `json:"announced"`
	TS                string     `json:"ts"`
	Name              string     `json:"name"`
	PlayOrder         int        `json:"play_order"`
	Level             string     `json:"level"`
	SetScores         []SetScore `json:"scores"`
}

type LiveGameData struct {
	Id             int    `json:"id"`
	CurrentSet     int    `json:"current_set"`
	HomePlayerName string `json:"home_player_name"`
	AwayPlayerName string `json:"away_player_name"`
	Phase          string `json:"phase"`
	GroupName      string `json:"group_name"`
}

type SetScore struct {
	Set  int `json:"set"`
	Home int `json:"home"`
	Away int `json:"away"`
}

type SetFinal struct {
	GameId       int `json:"game_id"`
	WinsRequired int `json:"wins_required"`
	Home         int `json:"home"`
	Away         int `json:"away"`
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
	S6hp *int `json:"s6hp"`
	S6ap *int `json:"s6ap"`
	S7hp *int `json:"s7hp"`
	S7ap *int `json:"s7ap"`
}

type GameMode struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ShortName    string `json:"short_name"`
	WinsRequired int    `json:"wins_required"`
	MaxSets      int    `json:"max_sets"`
}

type GameTimeline struct {
	Summary GameTimelineGameSummary `json:"summary"`
	Sets    map[int]*Set            `json:"sets"`
}

type Set struct {
	Events     []*GameEvent `json:"events"`
	SetSummary SetSummary   `json:"set_summary"`
}

type SetSummary struct {
	HomePoints          int    `json:"home_points"`
	AwayPoints          int    `json:"away_points"`
	HomeServes          int    `json:"home_serves"`
	AwayServes          int    `json:"away_serves"`
	HomeServePoints     int    `json:"home_serve_points"`
	AwayServePoints     int    `json:"away_serve_points"`
	HomeStreak          int    `json:"home_streak"`
	AwayStreak          int    `json:"away_streak"`
	StartTimestamp      int    `json:"start_timestamp"`
	EndTimestamp        int    `json:"end_timestamp"`
	HomeServePointsPerc string `json:"home_serve_points_perc"`
	AwayServePointsPerc string `json:"away_serve_points_perc"`
}

type GameTimelineGameSummary struct {
	GameStartingServerId    int    `json:"game_starting_server_id"`
	WinnerId                int    `json:"winner_id"`
	HomePlayerId            int    `json:"home_player_id"`
	AwayPlayerId            int    `json:"away_player_id"`
	HomeName                string `json:"home_name"`
	AwayName                string `json:"away_name"`
	GroupName               string `json:"group_name"`
	TournamentName          string `json:"tournament_name"`
	HomeTotalScore          int    `json:"home_total_score"`
	AwayTotalScore          int    `json:"away_total_score"`
	HomeTotalPoints         int    `json:"home_total_points"`
	AwayTotalPoints         int    `json:"away_total_points"`
	HomePointsPerc          string `json:"home_points_perc"`
	AwayPointsPerc          string `json:"away_points_perc"`
	HomeServesTotal         int    `json:"home_serves_total"`
	AwayServesTotal         int    `json:"away_serves_total"`
	HomeOwnServePointsTotal int    `json:"home_own_serve_points_total"`
	AwayOwnServePointsTotal int    `json:"away_own_serve_points_total"`
	HomeServePointsPerc     string `json:"home_serve_points_perc"`
	AwayServePointsPerc     string `json:"away_serve_points_perc"`
}

type GameEvent struct {
	CurrentSetStartingServer int `json:"current_set_starting_server"`
	CurrentServer            int `json:"current_server"`
	IsHomePoint              int `json:"is_home_point"`
	IsAwayPoint              int `json:"is_away_point"`
	HomePointsScored         int `json:"home_points_scored"`
	AwayPointsScored         int `json:"away_points_scored"`
	HomePlayerId             int `json:"home_player_id"`
	AwayPlayerId             int `json:"away_player_id"`
	Timestamp                int `json:"timestamp"`
	RallySeconds             int `json:"rally_seconds"`
	SetNumber                int `json:"set_number"`
}

type GameId struct {
	Id int `json:"id"`
}

type GameSetResults struct {
	GameId     int  `json:"game_id"`
	IsWalkover int  `json:"is_walkover"`
	S1hp       *int `json:"s1hp"`
	S1ap       *int `json:"s1ap"`
	S2hp       *int `json:"s2hp"`
	S2ap       *int `json:"s2ap"`
	S3hp       *int `json:"s3hp"`
	S3ap       *int `json:"s3ap"`
	S4hp       *int `json:"s4hp"`
	S4ap       *int `json:"s4ap"`
	S5hp       *int `json:"s5hp"`
	S5ap       *int `json:"s5ap"`
	S6hp       *int `json:"s6hp"`
	S6ap       *int `json:"s6ap"`
	S7hp       *int `json:"s7hp"`
	S7ap       *int `json:"s7ap"`
}

func (gsr GameSetResults) GetFullScore() (int, int) {
	var hs, as int
	if gsr.S1hp != nil && gsr.S1ap != nil {
		if *gsr.S1hp > *gsr.S1ap {
			hs++
		} else {
			as++
		}
	}
	if gsr.S2hp != nil && gsr.S2ap != nil {
		if *gsr.S2hp > *gsr.S2ap {
			hs++
		} else {
			as++
		}
	}
	if gsr.S3hp != nil && gsr.S3ap != nil {
		if *gsr.S3hp > *gsr.S3ap {
			hs++
		} else {
			as++
		}
	}
	if gsr.S4hp != nil && gsr.S4ap != nil {
		if *gsr.S4hp > *gsr.S4ap {
			hs++
		} else {
			as++
		}
	}
	if gsr.S5hp != nil && gsr.S5ap != nil {
		if *gsr.S5hp > *gsr.S5ap {
			hs++
		} else {
			as++
		}
	}
	if gsr.S6hp != nil && gsr.S6ap != nil {
		if *gsr.S6hp > *gsr.S6ap {
			hs++
		} else {
			as++
		}
	}
	if gsr.S7hp != nil && gsr.S7ap != nil {
		if *gsr.S7hp > *gsr.S7ap {
			hs++
		} else {
			as++
		}
	}

	return hs, as
}

type ChangeServerPayload struct {
	GameId   int `json:"game_id"`
	ServerId int `json:"server_id"`
}

type EloCache struct {
	Id             int  `json:"id"`
	HomePlayerId   int  `json:"home_player_id"`
	AwayPlayerId   int  `json:"away_player_id"`
	WinnerId       int  `json:"winner_id"`
	HomeScoreTotal int  `json:"home_score_total"`
	AwayScoreTotal int  `json:"away_score_total"`
	HomeElo        *int `json:"home_elo"`
	AwayElo        *int `json:"away_elo"`
	NewHomeElo     *int `json:"new_home_elo"`
	NewAwayElo     *int `json:"new_away_elo"`
	GamesPlayed    int  `json:"games_played"`
}

type EloHistory struct {
	GameId         int `json:"game_id"`
	HomePlayerId   int `json:"home_player_id"`
	AwayPlayerId   int `json:"away_player_id"`
	HomeScoreTotal int `json:"home_score_total"`
	AwayScoreTotal int `json:"away_score_total"`
	HomeEloOld     int `json:"home_elo"`
	AwayEloOld     int `json:"away_elo"`
	HomeEloNew     int `json:"new_home_elo"`
	AwayEloNew     int `json:"new_away_elo"`
	HomePlayed     int `json:"home_played"`
	AwayPlayed     int `json:"away_played"`
}

type ServeData struct {
	GameId                 int `json:"game_id"`
	SetNumber              int `json:"set_number"`
	FirstGameServer        int `json:"first_game_server"`
	SecondGameServer       int `json:"second_game_server"`
	CurrentSetFirstServer  int `json:"current_set_first_server"`
	CurrentSetSecondServer int `json:"current_set_second_server"`
	CurrentServerId        int `json:"current_server_id"`
	NumServes              int `json:"num_serves"`
}
