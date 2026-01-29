package models

type Tournament struct {
	Id                      int     `json:"id"`
	Name                    string  `json:"name"`
	StartTime               string  `json:"start_time"`
	IsPlayoffs              int     `json:"is_playoffs"`
	OfficeId                int     `json:"office_id"`
	Phase                   string  `json:"phase"`
	IsFinished              int     `json:"is_finished"`
	ParentTournamentId      *int    `json:"parent_tournament_id"`
	Participants            int     `json:"participants"`
	Scheduled               int     `json:"scheduled"`
	Finished                int     `json:"finished"`
	Sets                    int     `json:"sets"`
	Points                  int     `json:"points"`
	EnableTimelinessBonus   int     `json:"enable_timeliness_bonus"`
	TimelinessBonusEarly    float64 `json:"timeliness_bonus_early"`
	TimelinessBonusOntime   float64 `json:"timeliness_bonus_ontime"`
	TimelinessWindowDays    int     `json:"timeliness_window_days"`
}

type TournamentGroupSchedule struct {
	Name         string         `json:"name"`
	GameSchedule []GameSchedule `json:"game_schedule"`
}

type TournamentPerformanceGroup struct {
	Id              int                 `json:"group_id"`
	Name            string              `json:"group_name"`
	Abbreviation    string              `json:"group_abbreviation"`
	GroupPromotions int                 `json:"group_promotions"`
	GroupAvgElo     int                 `json:"group_avg_elo"`
	Players         []PlayerPerformance `json:"players"`
}

type PlayerPerformance struct {
	Pos               int    `json:"pos"`
	PlayerId          int    `json:"player_id"`
	PlayerName        string `json:"player_name"`
	StartingElo       int    `json:"starting_elo"`
	LastElo           int    `json:"last_elo"`
	GroupId           int    `json:"group_id"`
	GroupName         string `json:"group_name"`
	GroupAbbreviation string `json:"group_abbreviation"`
	Won               int    `json:"won"`
	Draw              int    `json:"draw"`
	Lost              int    `json:"lost"`
	Finished          int    `json:"finished"`
	Unfinished        int    `json:"unfinished"`
	Performance       int    `json:"performance"`
	Points            int    `json:"points"`
	TotalPoints       int    `json:"total_points"`
	Form              []int  `json:"form"`
}

type TournamentGroup struct {
	Id              int                    `json:"group_id"`
	Name            string                 `json:"group_name"`
	Abbreviation    string                 `json:"group_abbreviation"`
	GroupPromotions int                    `json:"group_promotions"`
	Players         []GroupStandingsPlayer `json:"players"`
}

type GroupStandingsPlayer struct {
	Pos                int     `json:"pos"`
	PosColor           string  `json:"pos_color"`
	PlayerId           int     `json:"player_id"`
	PlayerName         string  `json:"player_name"`
	Played             int     `json:"played"`
	Wins               int     `json:"wins"`
	Draws              int     `json:"draws"`
	Losses             int     `json:"losses"`
	Points             float64 `json:"points"`
	SetsFor            int     `json:"sets_for"`
	SetsAgainst        int     `json:"sets_against"`
	SetsDiff           int     `json:"sets_diff"`
	RalliesFor         *int    `json:"rallies_for"`
	RalliesAgainst     *int    `json:"rallies_against"`
	RalliesDiff        *int    `json:"rallies_diff"`
	GroupId            int     `json:"group_id"`
	GroupName          string  `json:"group_name"`
	GroupAbbreviation  string  `json:"group_abbreviation"`
	GroupPromotions    int     `json:"group_promotions"`
	GamesRemaining     int     `json:"games_remaining"`
	PointsPotentialMin float64 `json:"points_potential_min"`
	PointsPotentialMax float64 `json:"points_potential_max"`
	PositionMin        int     `json:"position_min"`
	PositionMax        int     `json:"position_max"`
}

type Ladder struct {
	GroupId     int           `json:"group_id"`
	GroupName   string        `json:"group_name"`
	LadderGroup []LadderGroup `json:"ladder_group"`
}

type LadderGroup struct {
	Order          int        `json:"order"`
	GameId         int        `json:"game_id"`
	GameName       string     `json:"game_name"`
	MaxStage       int        `json:"max_stage"`
	Stage          int        `json:"stage"`
	AwayPlayerId   int        `json:"away_player_id"`
	HomePlayerId   int        `json:"home_player_id"`
	WinnerId       int        `json:"winner_id"`
	HomeScoreTotal int        `json:"home_score_total"`
	AwayScoreTotal int        `json:"away_score_total"`
	IsWalkover     int        `json:"is_walkover"`
	HomePlayerName string     `json:"home_player_name"`
	AwayPlayerName string     `json:"away_player_name"`
	Level          string     `json:"level"`
	GroupName      string     `json:"group_name"`
	Announced      int        `json:"announced"`
	IsFinalGame    int        `json:"is_final_game"`
	SetScores      []SetScore `json:"scores"`
}

type GroupInfo struct {
	Id                 int                        `json:"id"`
	Name               string                     `json:"name"`
	GroupPromotions    int                        `json:"group_promotions"`
	GamesPlayed        int                        `json:"games_played"`
	GamesRemaining     int                        `json:"games_remaining"`
	PositionsUp        int                        `json:"positions_up"`
	PositionsStay      int                        `json:"positions_stay"`
	PositionsDown      int                        `json:"positions_down"`
	PlayerInfo         []PlayerInfo               `json:"player_info"`
	TopDrop            int                        `json:"top_drop"`
	TopDropPlayerName  string                     `json:"top_drop_player_name"`
	TopClimb           int                        `json:"top_climb"`
	TopClimbPlayerName string                     `json:"top_climb_player_name"`
	StatsMessage       string                     `json:"stats_message"`
	CandidatesMessage  string                     `json:"candidates_message"`
	PositionCandidates map[int]*PositionCandidate `json:"position_candidates"`
}

type PlayerInfo struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	PositionPrevious int    `json:"position_previous"`
	PositionCurrent  int    `json:"position_current"`
	PositionMovement int    `json:"position_movement"`
	PositionMin      int    `json:"position_min"`
	PositionMax      int    `json:"position_max"`
}

type PositionCandidate struct {
	PlayerNames []string `json:"player_names"`
	Secured     int      `json:"secured"`
}

type TournamentStatistics struct {
	Id                int     `json:"id"`
	Name              string  `json:"name"`
	Divisions         int     `json:"divisions"`
	SetsPlayed        int     `json:"sets_played"`
	PointsScored      int     `json:"points_scored"`
	AvgPointsPerMatch float32 `json:"avg_points_per_match"`
}

type TournamentPlayerStatistics struct {
	MostPointsGid                   int      `json:"most_points_gid"`
	MostPointsInGame                int      `json:"most_points_in_game"`
	MostPointsInGamePlayerId        int      `json:"most_points_in_game_player_id"`
	MostPointsInGamePlayerName      string   `json:"most_points_in_game_player_name"`
	LeastPointsGid                  int      `json:"least_points_gid"`
	LeastPointsLostInGame           int      `json:"least_points_lost_in_game"`
	LeastPointsLostInGamePlayerId   int      `json:"least_points_lost_in_game_player_id"`
	LeastPointsLostInGamePlayerName string   `json:"least_points_lost_in_game_player_name"`
	MostEloGainGid                  int      `json:"most_elo_gid"`
	MostEloGain                     int      `json:"most_elo_gain"`
	MostEloGainPlayerId             int      `json:"most_elo_gain_player_id"`
	MostEloGainPlayerName           string   `json:"most_elo_gain_player_name"`
	MostEloLostGid                  int      `json:"most_elo_lost_gid"`
	MostEloLost                     int      `json:"most_elo_lost"`
	MostEloLostPlayerId             int      `json:"most_elo_lost_player_id"`
	MostEloLostPlayerName           string   `json:"most_elo_lost_player_name"`
	MostPointsGameGid               int      `json:"most_points_game_gid"`
	MostPointsGame                  int      `json:"most_points_game"`
	MostPointsGameHomeName          string   `json:"most_points_game_home_name"`
	MostPointsGameAwayName          string   `json:"most_points_game_away_name"`
	LeastPointsGameGid              int      `json:"least_points_game_gid"`
	LeastPointsGame                 int      `json:"least_points_game"`
	LeastPointsGameHomeName         string   `json:"least_points_game_home_name"`
	LeastPointsGameAwayName         string   `json:"least_points_game_away_name"`
	MaxSetStreak                    int      `json:"max_set_streak"`
	MaxSetStreakPlayers             []Player `json:"max_set_streak_players"`
}
