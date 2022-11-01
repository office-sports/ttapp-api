package models

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
	EloHistory        [][2]int `json:"elo_history"`
}
