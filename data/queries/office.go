package queries

func GetOfficesQuery() string {
	return `SELECT
			o.id, o.name, o.is_default, o.channel_id
			FROM office o
			ORDER BY o.id ASC`
}

func GetOfficeQuery() string {
	return `SELECT
			o.id, o.name, o.is_default, o.channel_id
			FROM office o
			WHERE o.id = ?`
}
