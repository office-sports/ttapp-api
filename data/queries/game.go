package queries

func UpdateServerQuery() string {
	return `UPDATE game SET server_id = ? WHERE id = ?`
}

func FinishGameQuery() string {
	return `update game g set winner_id = ?, current_set = 1, is_finished = 1, date_played = now() where g.id = ?`
}

func CheckAndFinishTournament() string {
	return `update tournament t
	join
	(select t.id tournament_id, (count(g.id) - sum(g.is_finished)) as c
	from tournament t
	left join game g on t.id = g.tournament_id
	group by t.id) t2 on t2.tournament_id = t.id
	set t.is_finished = 1
	where c = 0 and t.id = ?`
}

func UpdateNextPlayoffGameHomePlayer() string {
	return `update game g set home_player_id = ? 
              where playoff_home_player_id = ? and tournament_id = ?`
}

func UpdateNextPlayoffGameAwayPlayer() string {
	return `update game g set away_player_id = ? 
              where playoff_away_player_id = ? and tournament_id = ?`
}

func GetLiveGamesQuery() string {
	return `select
		g.id, g.current_set, p1.name as homePlayerName, p2.name as awayPlayerName,
		if(t.is_playoffs = 1, 'playoffs', 'group') as phase, tg.name
		from game g
			join tournament t on g.tournament_id = t.id
			join tournament_group tg on g.tournament_group_id = tg.id
			join player p1 on p1.id = g.home_player_id
			join player p2 on p2.id = g.away_player_id
			join scores s on s.game_id = g.id
	where g.is_finished = 0
	and g.is_abandoned = 0
	having count(s.id > 0)`
}

func GetTournamentLiveGamesQuery() string {
	return `select
		g.id, g.current_set, p1.name as homePlayerName, p2.name as awayPlayerName,
		if(t.is_playoffs = 1, 'playoffs', 'group') as phase, tg.name
		from game g
			join tournament t on g.tournament_id = t.id
			join tournament_group tg on g.tournament_group_id = tg.id
			join player p1 on p1.id = g.home_player_id
			join player p2 on p2.id = g.away_player_id
	where g.is_finished = 0
	and g.is_abandoned = 0
	and g.tournament_id = ?
	and g.announced = 1`
}

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
	return `select g.server_id, g.winner_id, g.home_player_id, g.away_player_id,
       p1.name homeName, p2.name as awayName, tg.name as groupName, t.name as tournamentName,
       g.home_score as homeTotalScore, 
       g.away_score as awayTotalScore,
       coalesce(sum(s.home_points), 0) as homeTotalPoints,
       coalesce(sum(s.away_points), 0) as awayTotalPoints,
       coalesce((sum(s.home_points) / (sum(s.home_points) + sum(s.away_points))) * 100, 0) as homePointsPerc,
       coalesce((sum(s.away_points) / (sum(s.home_points) + sum(s.away_points))) * 100, 0) as awayPointsPerc
		from game g
		join tournament_group tg on g.tournament_group_id = tg.id
		join tournament t on g.tournament_id = t.id
		left join scores s on s.game_id = g.id
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

func GetEloHistory() string {
	return `select g.id, g.home_player_id, g.away_player_id, g.home_score, g.away_score,
			coalesce(g.old_home_elo, 0), 
			coalesce(g.old_away_elo, 0), 
			coalesce(g.new_home_elo, 0), 
			coalesce(g.new_away_elo, 0)
			from game g
			join tournament t on t.id = g.tournament_id
			where 
			g.is_finished = 1 and is_abandoned = 0 and t.is_official = 1
			order by g.date_played`
}

func GetPlayersEloData() string {
	return `select
    coalesce(ROW_NUMBER() OVER(ORDER BY g.date_played), 0) AS games_played,
    coalesce(if (g.home_player_id = ?, new_home_elo, new_away_elo), 1500) as elo
	from game g
	join tournament t on t.id = g.tournament_id
	where
	  ? in (g.home_player_id, away_player_id)
	  and g.is_finished = 1
	  and is_abandoned = 0
	  and t.is_official = 1
	  and g.id != ?
	order by g.date_played desc
	limit 1`
}
