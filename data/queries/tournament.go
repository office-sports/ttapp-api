package queries

func GetBaseTournamentScheduleQuery() string {
	return `select
			g.tournament_id                                         as tournamentId,
			tt.office_id                                            as officeId,
			g.id as matchId,
			tg.name                                                 as groupName,
			DATE(g.date_of_match)                                   as dateOfMatch,
			p1.id                                                   as homePlayerId,
			p2.id                                                      awayPlayerId,
			p1.name                                                    homePlayerName,
			p2.name                                                 as awayPlayerName,
			g.home_score                                            as homeScoreTotal,
			g.away_score                                            as awayScoreTotal
			from game g
			join game_mode gm on gm.id = g.game_mode_id
			left join player p1 on p1.id = g.home_player_id
			left join player p2 on p2.id = g.away_player_id
			left join tournament_group tg on tg.id = g.tournament_group_id
			left join tournament tt on tt.id = g.tournament_id`
}

func GetTournamentStandingsQuery() string {
	return `select 
    			0 as pos, tg.color_template as colorTemplate, p.id as playerId, p.name as playerName, 
			    SUM(if (g1.is_finished = 1, 1, 0)) as played, 
				SUM(if (g1.winner_id = ptg.player_id, 1, 0)) as wins,
				SUM(if (g1.is_finished = 1 AND g1.winner_id = 0, 1, 0)) as draws,
				SUM(if (g1.is_finished = 1 AND g1.winner_id != 0 AND g1.winner_id != ptg.player_id, 1, 0)) as losses,
				(SUM(if (g1.winner_id = ptg.player_id, 1, 0)) * 2 + SUM(if (g1.is_finished = 1 AND g1.winner_id = 0, 1, 0))) as points,
				(SUM(if (g1.home_player_id = ptg.player_id, g1.home_score, 0)) + SUM(if (g1.away_player_id = ptg.player_id, g1.away_score, 0))) as setsFor,
				(SUM(if (g1.home_player_id = ptg.player_id, g1.away_score, 0)) + SUM(if (g1.away_player_id = ptg.player_id, g1.home_score, 0))) as setsAgainst,
				((SUM(if (g1.home_player_id = ptg.player_id, g1.home_score, 0)) + SUM(if (g1.away_player_id = ptg.player_id, g1.away_score, 0))) -
				(SUM(if (g1.home_player_id = ptg.player_id, g1.away_score, 0)) + SUM(if (g1.away_player_id = ptg.player_id, g1.home_score, 0)))) as setDf,
				u.ralliesFor, u.ralliesAgainst, u.df,		    
			    ptg.group_id as groupId, tg.name as groupName, tg.abbreviation as groupAbbreviation
			from player_tournament_group ptg
			left join game g1 on (g1.home_player_id = ptg.player_id or g1.away_player_id = ptg.player_id) and g1.tournament_id = ?
			left join player p on p.id = ptg.player_id
			left join tournament_group tg on tg.id = ptg.group_id
			left join (
			select player, sum(pointsFor) as ralliesFor, sum(pointsAgainst) as ralliesAgainst, (sum(pointsFor) - sum(pointsAgainst)) as df from (
			SELECT g.id, g.home_player_id   AS player, sum(s.home_points) AS pointsFor, sum(s.away_points) AS pointsAgainst
			FROM scores s JOIN game g ON g.id = s.game_id
			JOIN tournament t1 on t1.id = g.tournament_id WHERE t1.id = ?
			GROUP BY g.home_player_id
			UNION
			SELECT
			g.id, g.away_player_id   AS player, sum(s.away_points) AS pointsFor, sum(s.home_points) AS pointsAgainst
			FROM scores s JOIN game g ON g.id = s.game_id
			JOIN tournament t2 on t2.id = g.tournament_id WHERE t2.id = ?
			GROUP BY g.away_player_id
			) u group by player
			) u on u.player = ptg.player_id
			where ptg.tournament_id = ?
			group by ptg.player_id
			order by ptg.group_id asc, points desc, setDf desc, u.df desc`
}

func GetTournamentGroupsQuery() string {
	return `select id, name from tournament_group tg where tg.tournament_id = ?`
}

func GetTournamentGroupQuery() string {
	return `select g.play_order as matchNumber, g.id, g.name, 4 as maxStage,
			g.stage, g.home_player_id as hpid, g.away_player_id as apid, g.winner_id,
			g.home_score as homeScoreTotal, g.away_score as awayScoreTotal, g.is_walkover,
			if (g.home_player_id, p1.name, g.playoff_home_player_id) as homePlayerDisplayName, 
			if (g.away_player_id, p2.name, g.playoff_away_player_id) as awayPlayerDisplayName  
			from game g 
			left join player p1 on p1.id = g.home_player_id 
			left join player p2 on p2.id = g.away_player_id 
			join tournament t on t.id = g.tournament_id and t.is_playoffs = 1
			join tournament_group l on l.id = g.tournament_group_id 
			where t.id = ? and l.id = ?
			order by g.stage asc, g.play_order desc`
}

func GetTournamentResultsQuery() string {
	return `SELECT
			g.id, tg.name as groupName, g.tournament_id as tournamentId, g.office_id as officeId,
			DATE(g.date_of_match) as dateOfMatch, DATE(g.date_played) as datePlayed,
			p1.id as homePlayerId, p2.id awayPlayerId, p1.name homePlayerName, p2.name as awayPlayerName,
			g.winner_id as winnerId, g.home_score as homeScoreTotal, g.away_score as awayScoreTotal,
			g.is_walkover as isWalkover,
			g.old_home_elo ohElo, g.new_home_elo nhElo, g.old_away_elo oaElo, g.new_away_elo naElo,
			(g.new_home_elo - g.old_home_elo) as homeEloDiff,
			(g.new_away_elo - g.old_away_elo) as awayEloDiff,
			if(count(ppp.id) > 0, 1, 0) as hasPoints,
			s1.home_points as s1hp, s1.away_points s1ap,
			s2.home_points as s2hp, s2.away_points s2ap,
			s3.home_points as s3hp, s3.away_points s3ap,
			s4.home_points as s4hp, s4.away_points s4ap,
			s5.home_points as s5hp, s5.away_points s5ap,
			s6.home_points as s6hp, s6.away_points s6ap,
			s7.home_points as s7hp, s7.away_points s7ap
			from game g
			join tournament t on g.tournament_id = t.id
			join game_mode gm on gm.id = g.game_mode_id
			left join player p1 on p1.id = g.home_player_id
			left join player p2 on p2.id = g.away_player_id
			left join tournament_group tg on tg.id = g.tournament_group_id
			left join scores s1 on s1.game_id = g.id and s1.set_number = 1
			left join scores s2 on s2.game_id = g.id and s2.set_number = 2
			left join scores s3 on s3.game_id = g.id and s3.set_number = 3
			left join scores s4 on s4.game_id = g.id and s4.set_number = 4
			left join scores s5 on s5.game_id = g.id and s5.set_number = 5
			left join scores s6 on s6.game_id = g.id and s6.set_number = 6
			left join scores s7 on s7.game_id = g.id and s7.set_number = 7
			left join scores sss on sss.game_id = g.id
			left join points ppp on ppp.score_id = sss.id`
}
