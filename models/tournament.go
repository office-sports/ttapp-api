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
