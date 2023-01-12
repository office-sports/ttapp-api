package queries

func GetOfficesQuery() string {
	return `SELECT
			o.id, o.name, o.is_default
			FROM office o
			ORDER BY o.id ASC`
}
