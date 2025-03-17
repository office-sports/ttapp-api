package models

type Leader struct {
	PlayerId      int    `json:"player_id"`
	PlayerName    string `json:"player_name"`
	ProfilePicUrl string `json:"profile_pic_url"`
	OfficeId      int    `json:"office_id"`
	GWon          int    `json:"g_won"`
	GLost         int    `json:"g_lost"`
	GDiff         int    `json:"g_diff"`
	PWon          int    `json:"p_won"`
	PLost         int    `json:"p_lost"`
	PDiff         int    `json:"p_diff"`
	SWon          int    `json:"s_won"`
	SLost         int    `json:"s_lost"`
	SDiff         int    `json:"s_diff"`
}

type Player struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Nickname          *string  `json:"nickname"`
	Elo               int      `json:"elo"`
	OldElo            int      `json:"old_elo"`
	EloChange         int      `json:"elo_change"`
	GamesPlayed       int      `json:"games_played"`
	Wins              int      `json:"wins"`
	Draws             int      `json:"draws"`
	Losses            int      `json:"losses"`
	OfficeId          int      `json:"office_id,omitempty"`
	WinPercentage     float32  `json:"win_percentage"`
	ProfilePicUrl     string   `json:"profile_pic_url"`
	NotWinPercentage  float32  `json:"not_win_percentage"`
	DrawPercentage    float32  `json:"draw_percentage"`
	NotDrawPercentage float32  `json:"not_draw_percentage"`
	LossPercentage    float32  `json:"loss_percentage"`
	NotLossPercentage float32  `json:"not_loss_percentage"`
	PPS               float32  `json:"pps"`
	EloHistory        [][2]int `json:"elo_history"`
	Active            int      `json:"active"`
}

type PlayerAvailability struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	GameDates []GameDate `json:"dates"`
}

type GameDate struct {
	Date string `json:"date"`
}

type PlayerDate struct {
	PlayerId int    `json:"player_id"`
	Date     string `json:"date"`
}

type PlayerCache struct {
	Elo         int `json:"elo"`
	OldElo      int `json:"old_elo"`
	GamesPlayed int `json:"games_played"`
}

type PlayerOpponent struct {
	Id                int    `json:"id"`
	Games             int    `json:"games"`
	OpponentId        int    `json:"opponent_id"`
	OpponentName      string `json:"opponent_name"`
	Wins              int    `json:"wins"`
	Draws             int    `json:"draws"`
	Losses            int    `json:"losses"`
	PlayerWalkovers   int    `json:"player_walkovers"`
	OpponentWalkovers int    `json:"opponent_walkovers"`
}
