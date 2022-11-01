package data

import "github.com/office-sports/ttapp-api/models"

func GetOffices() ([]*models.Office, error) {
	rows, err := models.DB.Query(`
			SELECT
			o.id, o.name, o.is_default
			FROM office o
			ORDER BY o.id ASC`)

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
