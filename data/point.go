package data

import (
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
	"log"
)

func InsertPoint(sid int, hp int, ap int) (int64, error) {
	s := `INSERT INTO points (score_id, is_home_point, is_away_point, time) VALUES (?, ?, ?, NOW())`
	lid, err := RunTransaction(s, sid, hp, ap)

	return lid, err
}

func GetMaxPoint(gid int, hp int, ap int) (int, error) {
	var pid int
	err := models.DB.QueryRow(queries.GetMaxPointQuery(), hp, ap, gid).Scan(&pid)

	if err != nil {
		return 0, err
	}

	return pid, nil
}

func DeletePointById(id int) {
	s := `DELETE FROM points WHERE id = ?`
	args := make([]interface{}, 0)
	args = append(args, id)
	_, err := RunTransaction(s, args...)

	if err != nil {
		log.Println("Error deleting point id: ", id, err)
	}
}
