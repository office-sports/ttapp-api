package models

type Tournament struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	StartTime    string `json:"start_time"`
	IsPlayoffs   int    `json:"is_playoffs"`
	OfficeId     int    `json:"office_id"`
	Phase        string `json:"phase"`
	IsFinished   int    `json:"is_finished"`
	Participants int    `json:"participants"`
	Scheduled    int    `json:"scheduled"`
	Finished     int    `json:"finished"`
	Sets         int    `json:"sets"`
	Points       int    `json:"points"`
}

type TournamentGroup struct {
	Id           int                    `json:"group_id"`
	Name         string                 `json:"group_name"`
	Abbreviation string                 `json:"group_abbreviation"`
	Players      []GroupStandingsPlayer `json:"players"`
}

type GroupStandingsPlayer struct {
	Pos               int    `json:"pos"`
	PosColor          string `json:"pos_color"`
	PlayerId          int    `json:"player_id"`
	PlayerName        string `json:"player_name"`
	Played            int    `json:"played"`
	Wins              int    `json:"wins"`
	Draws             int    `json:"draws"`
	Losses            int    `json:"losses"`
	Points            int    `json:"points"`
	SetsFor           int    `json:"sets_for"`
	SetsAgainst       int    `json:"sets_against"`
	SetsDiff          int    `json:"sets_diff"`
	RalliesFor        *int   `json:"rallies_for"`
	RalliesAgainst    *int   `json:"rallies_against"`
	RalliesDiff       *int   `json:"rallies_diff"`
	GroupId           int    `json:"group_id"`
	GroupName         string `json:"group_name"`
	GroupAbbreviation string `json:"group_abbreviation"`
}

type Ladder struct {
	GroupId     int           `json:"group_id"`
	GroupName   string        `json:"group_name"`
	LadderGroup []LadderGroup `json:"ladder_group"`
}

type LadderGroup struct {
	Order                 int    `json:"order"`
	GameId                int    `json:"game_id"`
	GameName              string `json:"game_name"`
	MaxStage              int    `json:"max_stage"`
	Stage                 int    `json:"stage"`
	AwayPlayerId          int    `json:"away_player_id"`
	HomePlayerId          int    `json:"home_player_id"`
	WinnerId              int    `json:"winner_id"`
	HomeScoreTotal        int    `json:"home_score_total"`
	AwayScoreTotal        int    `json:"away_score_total"`
	IsWalkover            int    `json:"is_walkover"`
	HomePlayerDisplayName string `json:"home_player_display_name"`
	AwayPlayerDisplayName string `json:"away_player_display_name"`
}

type LeaderGroup struct {
	Name      string   `json:"name"`
	LeaderSet []Leader `json:"leader_set"`
}

type Leader struct {
	PlayerId int `json:"player_id"`
	Value    int `json:"value"`
}
