package queries

func GetGameModesQuery() string {
	return `SELECT
			t.id, t.name, t.short_name, t.wins_required, t.max_sets
			FROM game_mode t
			ORDER BY t.id ASC`
}

func GetGameScoresQuery() string {
	return `SELECT id, game_id, set_number, home_points, away_points from scores s `
}

func GetGameServeDataQuery() string {
	return `select g.id, if(g.current_set = 0, 1, g.current_set) as setNumber,
                       @firstServer:=if(g.server_id = g.home_player_id, g.home_player_id, g.away_player_id) as firstGameServer,
                       @otherServer:=if(g.server_id = g.home_player_id, g.away_player_id, g.home_player_id) as secondGameServer,
                       @css:=if (mod(g.current_set + 1, 2) = 0, @firstServer, @otherServer) as currentSetFirstServer,
                       @oss:=if (mod(g.current_set + 1, 2) = 0, @otherServer, @firstServer) as currentSetSecondServer,
                       if (if (count(p.id) >= 20, 1, 0) = 1,
                           if(
                               mod(count(p.id), 2) = 0,
                               if (mod(g.current_set + 1, 2) = 0, @firstServer, @otherServer),
                               if (mod(g.current_set + 1, 2) = 0, @otherServer, @firstServer)
                               ),
                           if(mod((floor(count(p.id) / 2)), 2) = 0,
                               if (mod(g.current_set + 1, 2) = 0, @firstServer, @otherServer),
                               if (mod(g.current_set + 1, 2) = 0, @otherServer, @firstServer)
                               )
                           ) as currentServerId,
                       if (if(count(p.id) >= 20, 1, 0) = 1, 1, mod(count(p.id)+1, 2) + 1) as numServes
                from game g
                left join scores s on g.id = s.game_id and s.set_number = g.current_set
                left join points p on s.id = p.score_id
                where g.id = ?`
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

func GetEloCache() string {
	return `select
                g.id,
                home_player_id,
                away_player_id,
                winner_id,
                home_score,
                away_score,
                old_home_elo,
                old_away_elo,
                new_home_elo,
                new_away_elo
            from game g join tournament t on g.tournament_id = t.id and t.is_official = 1
            where g.is_finished = 1
            order by g.tournament_id, g.date_played, g.date_of_match, g.id asc`
}

func GetEloLastCache() string {
	return `select
                g.id,
                home_player_id,
                away_player_id,
                winner_id,
                home_score,
                away_score,
                old_home_elo,
                old_away_elo,
                new_home_elo,
                new_away_elo,
                coalesce(sum(if(g.home_player_id = ?, 1, 0)) + sum(if(g.away_player_id = ?, 1, 0)), 0) as gamesPlayed
            from game g join tournament t on g.tournament_id = t.id and t.is_official = 1
            where g.is_finished = 1
            and ? IN (g.home_player_id, g.away_player_id)
            and g.id != ?
            order by g.date_played desc limit 1`
}
