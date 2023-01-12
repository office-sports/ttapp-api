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
