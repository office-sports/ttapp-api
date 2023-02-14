package data

import (
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
)

func GetScoresByGameId(gid int) ([]*models.Score, error) {
	rows, err := models.DB.Query(queries.GetGameScoresQuery()+
		` where s.game_id = ?`, gid)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scores := make([]*models.Score, 0)
	for rows.Next() {
		g := new(models.Score)
		err := rows.Scan(&g.Id, &g.GameId, &g.SetNumber, &g.HomeScore, &g.AwayScore)
		if err != nil {
			return nil, err
		}

		scores = append(scores, g)
	}

	if err != nil {
		return nil, err
	}

	return scores, nil
}
