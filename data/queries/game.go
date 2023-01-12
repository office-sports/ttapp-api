package queries

func GetGameModesQuery() string {
	return `SELECT
			t.id, t.name, t.short_name, t.wins_required, t.max_sets
			FROM game_mode t
			ORDER BY t.id ASC`
}
