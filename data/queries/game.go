package queries

func GetGameModesQuery() string {
	return `SELECT
			t.id, t.name, t.short_name, t.wins_required, t.max_sets
			FROM game_mode t
			ORDER BY t.id ASC`
}

func GetGameTimelineSummaryQuery() string {
	return `select
    g.server_id as serverId,
    p1.name homeName, p2.name as awayName, tg.name as groupName, t.name as tournamentName,
    g.home_score as homeTotalScore, g.away_score as awayTotalScore,
    sum(p.is_home_point) as homeTotalPoints, sum(p.is_away_point) as awayTotalPoints,
    (sum(p.is_home_point) / (sum(p.is_home_point) + sum(p.is_away_point))) * 100 as homePointsPerc,
    (sum(p.is_away_point) / (sum(p.is_home_point) + sum(p.is_away_point))) * 100 as awayPointsPerc
    from points p
    join scores s on p.score_id = s.id
    join game g on s.game_id = g.id
    join tournament_group tg on g.tournament_group_id = tg.id
    join tournament t on g.tournament_id = t.id
    join player p1 on g.home_player_id = p1.id
    join player p2 on g.away_player_id = p2.id`
}

func GetGameTimelineQuery() string {
	return `select
		p.is_home_point as isHomePoint, p.is_away_point as isAwayPoint,
    	g.home_player_id as homePlayerId, g.away_player_id as awayPlayerId,
		unix_timestamp(p.time) as timestamp,
		s.set_number as setNumber
		from points p
		join scores s on p.score_id = s.id
		join game g on s.game_id = g.id
		join tournament_group tg on g.tournament_group_id = tg.id
		join tournament t on g.tournament_id = t.id
		join player p1 on g.home_player_id = p1.id
		join player p2 on g.away_player_id = p2.id`
}
