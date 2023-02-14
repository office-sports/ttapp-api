package queries

func GetMaxPointQuery() string {
	return `select max(p.id) as pointId
            from points p
            join scores s on s.id = p.score_id
            where p.is_home_point = ? and p.is_away_point = ? and s.game_id = ?`
}
