package queries

func GetPlayersDataQuery() string {
	return `SELECT 
				p.id, p.name, p.nickname, p.tournament_elo as elo, p.tournament_elo_previous as oldElo, 
				(p.tournament_elo - p.tournament_elo_previous) as eloChange, 
				if (count(g.id) is null, 0, count(g.id)) as gamesPlayed,
				sum(if(g.winner_id = p.id, 1, 0)) as wins,
				sum(if(g.winner_id = 0, 1, 0)) as draws, 
				sum(if(g.winner_id != 0 and g.winner_id != p.id, 1, 0)) as losses,
				p.office_id as officeId,
				IF (count(g.id) > 0, (sum(if(g.winner_id = p.id, 1, 0)) / count(g.id)) * 100 , 0) as winPercentage,
				coalesce(sum(if(g.home_player_id = p.id, s.home_points, s.away_points)) / count(s.id), 0) as pps,
				p.active
			from player p
			left join game g on p.id in (g.home_player_id, g.away_player_id) and g.is_finished = 1
			left join tournament t on t.id = g.tournament_id and t.is_official = 1
			left join scores s on s.game_id = g.id
			group by p.id, p.name
			order by p.name`
}

func GetPlayerDataQuery() string {
	return `SELECT 
				p.id, p.name, count(g.id) as played, p.profile_pic_url as pic, p.tournament_elo as elo, 
				sum(if(p.id = g.winner_id, 1, 0)) as wins, 
				sum(if(g.winner_id = 0, 1, 0)) as draws, 
				sum(if(g.winner_id != 0 and g.winner_id != p.id, 1, 0)) as losses,
				IF (count(g.id) > 0, (sum(if(g.winner_id = p.id, 1, 0)) / count(g.id)) * 100 , 0) as winPercentage,
				coalesce(sum(if(g.home_player_id = p.id, s.home_points, s.away_points)) / count(s.id), 0) as pps
			from player p 
			left join game g on p.id in (g.home_player_id, g.away_player_id) and g.is_finished = 1 
			left join tournament t on g.tournament_id = t.id and t.is_official = 1 
			left join scores s on s.game_id = g.id
			where p.id = ?  
			group by p.id`
}

func GetPlayerOpponentsQuery() string {
	return `select
				count(g.id) as games,
				if (g.home_player_id = ?, g.away_player_id, g.home_player_id) as opponent,
				if (g.home_player_id = ?, pa.name, ph.name) as opponentName,
				sum(if (g.winner_id = ?, 1, 0)) as wins,
				sum(if (g.winner_id = 0, 1, 0)) as draws,
				sum(if (g.winner_id != ? and g.winner_id != 0, 1, 0)) as losses,
    			sum(if (g.winner_id != ? and g.is_walkover = 1, 1, 0)) as playerWalkovers,
				sum(if (g.winner_id = ? and g.is_walkover = 1, 1, 0)) as opponentWalkovers
			from game g
				join tournament t on t.id = g.tournament_id
			join player ph on ph.id = g.home_player_id
			join player pa on pa.id = g.away_player_id
			where ? in (g.home_player_id, g.away_player_id)
			and t.is_official = 1
			and g.is_finished = 1
			group by opponent
			order by count(g.id) desc`
}

func GetPlayerEloHistoryQuery() string {
	return `select 
    		coalesce(if (p.id = g.home_player_id, g.new_home_elo, g.new_away_elo), 0) as elo
			from player p
			left join game g on p.id in (g.home_player_id, g.away_player_id)
			join tournament t on g.tournament_id = t.id
			where p.id = ?
			and g.is_finished = 1
			and t.is_official = 1
			and g.is_abandoned = 0
			order by g.date_played asc, g.id asc`
}

func GetGameWithScoresQuery() string {
	return `select 
    	g.id, 
    	gm.max_sets,
    	gm.wins_required as winsRequired,
    	g.tournament_id as tournamentId,
    	g.office_id as officeId,
		tg.name as groupName, 
		g.date_of_match as dateOfMatch, 
		g.date_played as datePlayed,
		coalesce(p1.id, 0) as homePlayerId,
		coalesce(p2.id, 0) as  awayPlayerId,
    	coalesce(p1.name, g.playoff_home_player_id) homePlayerName,
    	coalesce(p2.name, g.playoff_away_player_id) awayPlayerName,
    	g.winner_id as winnerId, 
		g.home_score as homeScoreTotal, 
		g.away_score as awayScoreTotal,
		g.is_walkover as isWalkover,
		g.is_finished as isFinished,
		COALESCE(g.old_home_elo, p1.tournament_elo) as homeElo, 
		COALESCE(g.old_away_elo, p2.tournament_elo) as awayElo,
		g.new_home_elo as newHomeElo, 
		g.new_away_elo as newAwayElo,
		g.new_home_elo - g.old_home_elo as homeEloDiff,
		g.new_away_elo - g.old_away_elo as awayEloDiff,
		s1.home_points as s1hp, s1.away_points s1ap,
		s2.home_points as s2hp, s2.away_points s2ap,
		s3.home_points as s3hp, s3.away_points s3ap,
		s4.home_points as s4hp, s4.away_points s4ap,
		s5.home_points as s5hp, s5.away_points s5ap,
		s6.home_points as s5hp, s6.away_points s6ap,
		s7.home_points as s5hp, s7.away_points s7ap,
		coalesce(chp, 0) as currentHomePoints,
        coalesce(cap, 0) as currentAwayPoints,
		g.current_set as currentSet,
		s.id as currentSetId,
		if(count(ppp.id) > 0, 1, 0) as hasPoints,
		g.announced, g.ts, COALESCE(g.name, '') as name, COALESCE(g.play_order, 0) as play_order,
		coalesce(g.level, "") as level
		from game g
		join game_mode gm on gm.id = g.game_mode_id
		left join player p1 on p1.id = g.home_player_id
		left join player p2 on p2.id = g.away_player_id
		join tournament_group tg on tg.id = g.tournament_group_id
		left join scores s on s.set_number = g.current_set and s.game_id = g.id
		left join scores s1 on s1.game_id = g.id and s1.set_number = 1
		left join scores s2 on s2.game_id = g.id and s2.set_number = 2
		left join scores s3 on s3.game_id = g.id and s3.set_number = 3
		left join scores s4 on s4.game_id = g.id and s4.set_number = 4
		left join scores s5 on s5.game_id = g.id and s5.set_number = 5
		left join scores s6 on s6.game_id = g.id and s6.set_number = 6
		left join scores s7 on s7.game_id = g.id and s7.set_number = 7
		left join (
            select
                s.game_id as gid,
                s.set_number as sn,
                coalesce(sum(p.is_home_point), 0) as chp,
                coalesce(sum(p.is_away_point), 0) as cap
                from scores s
            left join points p on p.score_id = s.id
            group by s.game_id, s.set_number
		) ss on ss.gid = g.id and ss.sn = g.current_set
		left join points pp on pp.score_id = ss.gid
		left join scores sss on sss.game_id = g.id
		left join points ppp on ppp.score_id = sss.id`
}

func GetLeadersQuery() string {
	return `select playerId player_id, pp.name player_name, 
       		   pp.profile_pic_url, pp.office_id,
       		   sum(won) as g_won, sum(lost) as g_lost,
			   (sum(won) - sum(lost)) g_diff,
			   sum(pointsFor) p_won, sum(pointsAgainst) p_lost,
			   (sum(pointsFor) - sum(pointsAgainst)) p_diff,
			   sum(setsFor) sWon, sum(setsAgainst) s_lost,
			   (sum(setsFor) - sum(setsAgainst)) s_diff
			   from (
					SELECT g.id,
							  g.home_player_id   AS playerId,
							  sum(if (g.winner_id = g.home_player_id, 1, 0)) as won,
							  sum(if (g.winner_id = g.away_player_id, 1, 0)) as lost,
							  sum(g.home_score) as setsFor,
							  sum(g.away_score) as setsAgainst,
							  sum(s.home_points) AS pointsFor,
							  sum(s.away_points) AS pointsAgainst
					   FROM scores s
								JOIN game g ON g.id = s.game_id
								JOIN tournament t1 on t1.id = g.tournament_id
					   GROUP BY g.home_player_id
					   UNION
					   SELECT g.id,
							  g.away_player_id   AS playerId,
							  sum(if (g.winner_id = g.away_player_id, 1, 0)) as won,
							  sum(if (g.winner_id = g.home_player_id, 1, 0)) as lost,
							  sum(g.away_score) as setsFor,
							  sum(g.home_score) as setsAgainst,
							  sum(s.away_points) AS pointsFor,
							  sum(s.home_points) AS pointsAgainst
					   FROM scores s
								JOIN game g ON g.id = s.game_id
								JOIN tournament t2 on t2.id = g.tournament_id
					   GROUP BY g.away_player_id) p
		join player pp on pp.id = p.playerId
			   where pp.active = 1
		group by playerId
		having g_lost > 0`
}

func GetPlayerLastEloDataQuery() string {
	return `select
			ROW_NUMBER() OVER(ORDER BY g.date_played) AS games_played,
			g.id,
			g.home_player_id,
			g.away_player_id,
			g.home_score,
			g.away_score,
			g.old_home_elo,
			g.old_away_elo,
			g.new_home_elo,
			g.new_away_elo
			from game g
			join tournament t on t.id = g.tournament_id
			where
			  ? in (g.home_player_id, away_player_id)
			  and g.is_finished = 1
			  and is_abandoned = 0
			  and t.is_official = 1
			  and g.id != 1401
			order by g.date_played desc
			limit 1`
}
