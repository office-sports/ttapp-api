package data

import (
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
)

func GetGameModes() ([]*models.GameMode, error) {
	rows, err := models.DB.Query(queries.GetGameModesQuery())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gm := make([]*models.GameMode, 0)
	for rows.Next() {
		o := new(models.GameMode)
		err := rows.Scan(&o.ID, &o.Name, &o.ShortName, &o.WinsRequired, &o.MaxSets)
		if err != nil {
			return nil, err
		}

		gm = append(gm, o)
	}

	if err != nil {
		return nil, err
	}

	return gm, nil
}
