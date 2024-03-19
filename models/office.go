package models

type Office struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	IsDefault int     `json:"is_default"`
	ChannelId *string `json:"channel_id"`
}
