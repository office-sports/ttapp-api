package data

import (
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
)

func GetOffices() ([]*models.Office, error) {
	rows, err := models.DB.Query(queries.GetOfficesQuery())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offices := make([]*models.Office, 0)
	for rows.Next() {
		o := new(models.Office)
		err := rows.Scan(&o.ID, &o.Name, &o.IsDefault)
		if err != nil {
			return nil, err
		}

		offices = append(offices, o)
	}

	if err != nil {
		return nil, err
	}

	return offices, nil
}
