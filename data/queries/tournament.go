package queries

func GetBaseTournamentScheduleQuery() string {
	return `select
			g.tournament_id                                         as tournamentId,
			tt.office_id                                            as officeId,
			g.id as matchId,
			tg.name                                                 as groupName,
			g.date_of_match                                         as dateOfMatch,
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
