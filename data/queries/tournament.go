package queries

func GetTournamentGroupScheduleQuery() string {
	return `select tg.name, ph.name, pa.name, ph.slack_name, pa.slack_name,
			   (week(g.date_of_match) - start_week + 1) as game_week,
			   g.date_of_match, g.is_finished
			from game g
				join tournament t on t.id = g.tournament_id
				join tournament_group tg on tg.id = g.tournament_group_id
				join player ph on ph.id = g.home_player_id
				join player pa on pa.id = g.away_player_id
			join (
				select min(week(g.date_of_match)) as start_week, g.tournament_id
				from game g
				join tournament t on g.tournament_id = t.id
				group by g.tournament_id
			) w on w.tournament_id = g.tournament_id
				 where t.is_playoffs = 0 and t.is_finished = 0
				   and g.is_finished = 0
				 and g.office_id = ?
				 and week(g.date_of_match) <= week(NOW())
		order by tg.id, date_of_match`
}

func GetTournamentGroupGamesQuery() string {
	return `select tg.name, ph.id, pa.id, ph.name, pa.name, ph.slack_name, pa.slack_name,
			   (week(g.date_of_match) - start_week + 1) as game_week,
			   g.date_of_match, g.is_finished
			from game g
				join tournament t on t.id = g.tournament_id
				join tournament_group tg on tg.id = g.tournament_group_id
				join player ph on ph.id = g.home_player_id
				join player pa on pa.id = g.away_player_id
			join (
				select min(week(g.date_of_match)) as start_week, g.tournament_id
				from game g
				join tournament t on g.tournament_id = t.id
				group by g.tournament_id
			) w on w.tournament_id = g.tournament_id
				 where t.is_playoffs = 0 and t.is_finished = 0
				 and g.office_id = ?
		order by tg.id, date_of_match`
}

func GetStatsLongestSetStreak() string {
	return `select g.id, g.home_player_id, g.away_player_id, g.winner_id, DATE(g.date_played),
				   s.id, if (s.home_points > s.away_points, 1, 0) as homeWon, if (s.home_points < s.away_points, 1, 0) as awayWon
			from game g
			join scores s on s.game_id = g.id
			where g.tournament_id = ?
			and g.is_finished = 1
			order by g.date_played, s.set_number`
}

func GetStatsMostPointsGameQuery() string {
	return `select g.id, sum(s.home_points) + sum(s.away_points) as points,
				   hp.name, ap.name
			from game g
			join scores s on s.game_id = g.id
			join player hp on hp.id = g.home_player_id
			join player ap on ap.id = g.away_player_id
			where g.tournament_id = ?
			and g.is_finished = 1
			group by g.id
			order by points desc
			limit 1;`
}

func GetStatsLeastPointsGameQuery() string {
	return `select g.id, sum(s.home_points) + sum(s.away_points) as points,
				   hp.name, ap.name
			from game g
			join scores s on s.game_id = g.id
			join player hp on hp.id = g.home_player_id
			join player ap on ap.id = g.away_player_id
			where g.tournament_id = ?
			and g.is_finished = 1
			group by g.id
			order by points
			limit 1;`
}

func GetStatsMostPointsInGameQuery() string {
	return `select g.id, if (sum(s.home_points) > sum(s.away_points), sum(s.home_points), sum(s.away_points)) as pointsScored,
				   if(sum(s.home_points) > sum(s.away_points), hp.id, ap.id) as playerId,
				   if(sum(s.home_points) > sum(s.away_points), hp.name, ap.name) as playerName
			from game g
			join scores s on s.game_id = g.id
			join player hp on hp.id = g.home_player_id
			join player ap on ap.id = g.away_player_id
			where g.tournament_id = ?
			group by g.id
			order by pointsScored desc
			limit 1;`
}

func GetTournamentProbabilitiesQuery() string {
	return `select 
		g.id,
		g.is_finished,
		g.is_abandoned,
		CASE WHEN EXISTS (SELECT 1 FROM scores s WHERE s.game_id = g.id) THEN 1 ELSE 0 END as is_started,
		g.is_walkover,
		g.home_player_id,
		g.away_player_id,
		g.winner_id,
		coalesce(g.old_home_elo, 1500) as home_elo,
		coalesce(g.old_away_elo, 1500) as away_elo
	from game g
	where g.tournament_id = ?
	order by g.id`
}

func GetStatsEloGainQuery() string {
	return `select
				g.id, GREATEST((g.new_home_elo - g.old_home_elo), (g.new_away_elo - g.old_away_elo)) as eloGain,
				g.winner_id as playerId, p.name
			from game g
			join player p on p.id = g.winner_id
			where g.tournament_id = ?
			and g.is_finished = 1
			order by eloGain desc
			limit 1;`
}

func GetStatsEloLostQuery() string {
	return `select
				g.id, LEAST((g.new_home_elo - g.old_home_elo), (g.new_away_elo - g.old_away_elo)) as eloGain,
				if(g.winner_id = g.home_player_id, ap.id, hp.id) as playerId,
				if(g.winner_id = g.home_player_id, ap.name, hp.name) as playerName
			from game g
			join player hp on hp.id = g.home_player_id
			join player ap on ap.id = g.away_player_id
			where g.tournament_id = ?
			and g.is_finished = 1
			order by eloGain asc
			limit 1;`
}

func GetStatsLeastPointsLostInGameQuery() string {
	return `select g.id, sum(if(g.winner_id = g.home_player_id, s.away_points, s.home_points)) as pointsLost,
				   if(g.winner_id = g.home_player_id, hp.id, ap.id) as playerId,
				   if(g.winner_id = g.home_player_id, hp.name, ap.name) as playerName
			from game g
			join scores s on s.game_id = g.id
			join player hp on hp.id = g.home_player_id
			join player ap on ap.id = g.away_player_id
			where g.tournament_id = ?
			group by g.id
			order by pointsLost
			limit 1;`
}

func GetTournamentsStatisticsQuery() string {
	return `select t.id, t.name, 
				   tg.divisions, sum(g.home_score + g.away_score) as setsPlayed,
				   scores.points,
				   scores.points / sum(if(g.is_finished = 1, 1, 0)) as pointsPerMatch
			from tournament t
			join game g on g.tournament_id = t.id
			join (select tg.tournament_id, count(tg.id) as divisions
				  from tournament_group tg
				  group by tg.tournament_id
				  ) tg on tg.tournament_id = t.id
			join (select g.tournament_id, sum(s.home_points) + sum(s.away_points) as points
				  from scores s
				  join game g on g.id = s.game_id
				  group by g.tournament_id
				  ) scores on scores.tournament_id = g.tournament_id
			group by t.id;`
}

func GetTournamentsQuery() string {
	return `select t.id, t.name, t.start_time, t.is_playoffs, t.office_id,
			IF (t.is_playoffs = 0, 'group', 'playoffs') as phase,
			t.is_finished, t.parent_tournament, coalesce(participants, 0), count(g.id),
			if(sum(g.is_finished) is null, 0, sum(g.is_finished)), coalesce(s.sets, 0), coalesce(s.points, 0)
			from tournament t
			left join (
                select gg.tid, count(distinct(gg.pid)) participants from (select g.home_player_id pid, g.tournament_id tid
                                           from game g
                                           where g.home_player_id != 0
                                           UNION ALL
                                           select g.away_player_id pid, g.tournament_id tid
                                           from game g
                                           where g.away_player_id != 0) gg
                group by gg.tid
            ) t2 on t2.tid = t.id
			left join game g on g.tournament_id = t.id
			left join (
			    select g.tournament_id as tid, count(s.id) as sets, sum(s.home_points + s.away_points) as points
                from scores s
                join game g on s.game_id = g.id
                group by g.tournament_id
            ) s on s.tid = t.id
			where t.is_official = 1
			group by t.id, t.start_time
			order by t.start_time desc;`
}

func GetLiveTournamentsQuery() string {
	return `select t.id, t.name, t.is_finished, t.is_playoffs, t.start_time, 
	   IF (t.is_playoffs = 0, 'group', 'playoffs') as phase,
       t.office_id, t.parent_tournament,
			count(distinct (g.home_player_id)), count(g.id),
			if(sum(g.is_finished) is null, 0, sum(g.is_finished)), coalesce(s.sets, 0), coalesce(s.points, 0)
			from tournament t
			left join game g on g.tournament_id = t.id
			left join (
			    select g.tournament_id as tid, count(s.id) as sets, sum(s.home_points + s.away_points) as points
                from scores s
                join game g on s.game_id = g.id
                group by g.tournament_id
            ) s on s.tid = t.id				
			where t.is_official = 1 
			group by t.id
			order by t.is_finished asc, t.start_time asc`
}

func GetTournamentByIdQuery() string {
	return `
			select t.id, t.name, t.start_time, t.is_playoffs, t.office_id,
			IF (t.is_playoffs = 0, 'group', 'playoffs') as phase,
			t.is_finished, t.parent_tournament, coalesce(participants, 0), count(g.id),
			if(sum(g.is_finished) is null, 0, sum(g.is_finished)),
			t.enable_timeliness_bonus, t.timeliness_bonus_early, t.timeliness_bonus_ontime, t.timeliness_window_days
			from tournament t
			left join (
                select gg.tid, count(distinct(gg.pid)) participants from (select g.home_player_id pid, g.tournament_id tid
                                           from game g
                                           where g.home_player_id != 0
                                           UNION ALL
                                           select g.away_player_id pid, g.tournament_id tid
                                           from game g
                                           where g.away_player_id != 0) gg
                group by gg.tid
            ) t2 on t2.tid = t.id
			left join game g on g.tournament_id = t.id
			where t.is_official = 1 and t.id = ?
			group by t.id, t.start_time
			order by t.start_time desc`
}

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
			g.away_score                                            as awayScoreTotal,
			gm.max_sets as bo
			from game g
			join game_mode gm on gm.id = g.game_mode_id
			left join player p1 on p1.id = g.home_player_id
			left join player p2 on p2.id = g.away_player_id
			left join tournament_group tg on tg.id = g.tournament_group_id
			left join tournament tt on tt.id = g.tournament_id`
}

func GetBaseTournamentGroupScheduleQuery() string {
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
			g.away_score                                            as awayScoreTotal,
			gm.max_sets as bo
			from game g
			join game_mode gm on gm.id = g.game_mode_id
			left join player p1 on p1.id = g.home_player_id
			left join player p2 on p2.id = g.away_player_id
			left join tournament_group tg on tg.id = g.tournament_group_id
			left join tournament tt on tt.id = g.tournament_id
			where g.is_finished = 0 and tt.is_finished = 0
			and g.tournament_id = ?
			and tg.id = ?
			group by g.id order by g.date_of_match, g.id asc`
}

func GetTournamentStandingsBaseQuery() string {
	return `select 
    			0 as pos, tg.color_template as colorTemplate, p.id as playerId, p.name as playerName, 
			    SUM(if (g1.is_finished = 1, 1, 0)) as played, 
				SUM(if (g1.winner_id = ptg.player_id, 1, 0)) as wins,
				SUM(if (g1.is_finished = 1 AND g1.winner_id = 0, 1, 0)) as draws,
				SUM(if (g1.is_finished = 1 AND g1.winner_id != 0 AND g1.winner_id != ptg.player_id, 1, 0)) as losses,
				(
					SUM(if (g1.winner_id = ptg.player_id, 1, 0)) * 2 + 
					SUM(if (g1.is_finished = 1 AND g1.winner_id = 0, 1, 0)) +
					-- Timeliness bonus (only if enabled and not playoffs)
					CASE 
						WHEN t.enable_timeliness_bonus = 1 AND t.is_playoffs = 0 THEN
							SUM(
								CASE
									-- On-time: within window days (date-based, ignoring time)
									WHEN g1.is_finished = 1 AND 
										 ABS(TIMESTAMPDIFF(DAY, DATE(g1.date_played), DATE(g1.date_of_match))) <= t.timeliness_window_days 
									THEN t.timeliness_bonus_ontime
									
									-- Early: played before window starts (date-based)
									WHEN g1.is_finished = 1 AND 
										 TIMESTAMPDIFF(DAY, DATE(g1.date_played), DATE(g1.date_of_match)) > t.timeliness_window_days
									THEN t.timeliness_bonus_early
									
									-- Late or not finished: no bonus
									ELSE 0
								END
							)
						ELSE 0
					END
				) as points,
				(SUM(if (g1.home_player_id = ptg.player_id, g1.home_score, 0)) + SUM(if (g1.away_player_id = ptg.player_id, g1.away_score, 0))) as setsFor,
				(SUM(if (g1.home_player_id = ptg.player_id, g1.away_score, 0)) + SUM(if (g1.away_player_id = ptg.player_id, g1.home_score, 0))) as setsAgainst,
				((SUM(if (g1.home_player_id = ptg.player_id, g1.home_score, 0)) + SUM(if (g1.away_player_id = ptg.player_id, g1.away_score, 0))) -
				(SUM(if (g1.home_player_id = ptg.player_id, g1.away_score, 0)) + SUM(if (g1.away_player_id = ptg.player_id, g1.home_score, 0)))) as setDf,
				u.ralliesFor, u.ralliesAgainst, u.df,		    
			    ptg.group_id as groupId, tg.name as groupName, tg.abbreviation as groupAbbreviation,
			    tg.promotions
			from player_tournament_group ptg
			left join game g1 on (g1.home_player_id = ptg.player_id or g1.away_player_id = ptg.player_id) and g1.tournament_id = ?
			left join player p on p.id = ptg.player_id
			left join tournament_group tg on tg.id = ptg.group_id
			left join tournament t on t.id = ptg.tournament_id
			left join (
			select player, sum(pointsFor) as ralliesFor, sum(pointsAgainst) as ralliesAgainst, (sum(pointsFor) - sum(pointsAgainst)) as df from (
			SELECT g.id, g.home_player_id   AS player, sum(s.home_points) AS pointsFor, sum(s.away_points) AS pointsAgainst
			FROM scores s JOIN game g ON g.id = s.game_id
			JOIN tournament t1 on t1.id = g.tournament_id WHERE t1.id = ?
			GROUP BY g.home_player_id
			UNION
			SELECT
			g.id, g.away_player_id AS player, sum(s.away_points) AS pointsFor, sum(s.home_points) AS pointsAgainst
			FROM scores s JOIN game g ON g.id = s.game_id
			JOIN tournament t2 on t2.id = g.tournament_id WHERE t2.id = ?
			GROUP BY g.away_player_id
			) u group by player
			) u on u.player = ptg.player_id
			where ptg.tournament_id = ? `
}

func GetTournamentStandingsQuery() string {
	return GetTournamentStandingsBaseQuery() + `group by ptg.player_id
			order by ptg.group_id asc, points desc, setDf desc, u.df desc, p.id asc`
}

func GetPlayersTournamentEloQuery() string {
	return `
		select id, playerId, elo, lelo, winner_id from (
		select g.id,
							  g.home_player_id as playerId,
							  g.old_home_elo as elo,
							  g.new_home_elo as lelo,
							  g.date_played as dp,
							  g.winner_id
					   from game g
					   where g.tournament_id = ?
						 and g.is_finished = 1
					   union
					   select g.id,
							  g.away_player_id as playerId,
							  g.old_away_elo as elo,
							  g.new_away_elo as lelo,
							  g.date_played as dp,
							  g.winner_id
					   from game g
					   where g.tournament_id = ?
						 and g.is_finished = 1) gg
		order by gg.playerId, gg.dp
`
}

func GetTournamentPerformanceQuery() string {
	return `select
					0 pos, p.playerId, p2.name playerName, p2.current_elo,
					tg.id as group_id, tg.name groupName, tg.abbreviation,
				   sum(won) as won, (sum(finished - (won + lost))) as draw, sum(lost) as lost, sum(finished) as finished, sum(unfinished) as unfinished,
				   if (sum(finished) = 0, 0, round((sum(p.sumElo) + 400 * (sum(won) - sum(lost)))/sum(finished))) as performance,
				   (sum(won)*2) as points, (sum(finished) + sum(unfinished)) * 2 as totalPoints
			from (select g.tournament_group_id, g.home_player_id                                         playerId,
						 sum(if(g.is_finished = 1, old_away_elo, 0))           as sumElo,
						 sum(if(g.is_finished = 1, 1, 0))                      as finished,
						 sum(if(g.is_finished = 0, 1, 0))                      as unfinished,
						 sum(if(g.winner_id = g.home_player_id, 1, 0))                       as won,
						 sum(if(g.winner_id != g.home_player_id AND g.winner_id != 0, 1, 0)) as lost
				  from game g
				  where g.tournament_id = ?
				  group by g.home_player_id
				  union
				  select g.tournament_group_id, g.away_player_id                                         playerId,
						 sum(if(g.is_finished = 1, old_home_elo, 0))           as sumElo,
						 sum(if(g.is_finished = 1, 1, 0))                      as finished,
						 sum(if(g.is_finished = 0, 1, 0))                      as unfinished,
						 sum(if(g.winner_id = g.away_player_id, 1, 0))                       as won,
						 sum(if(g.winner_id != g.away_player_id AND g.winner_id != 0, 1, 0)) as lost
				  from game g
				  where g.tournament_id = ?
				  group by g.away_player_id) p
				join player p2 on p2.id = p.playerId
			left join player_tournament_group ptg on ptg.player_id = p.playerId and ptg.tournament_id = ?
			left join tournament_group tg on tg.id = coalesce(ptg.group_id, p.tournament_group_id)
			group by p.playerId
			order by ptg.group_id asc, performance desc`
}

func GetTournamentStandingsDaysQuery() string {
	return GetTournamentStandingsBaseQuery() +
		`and (g1.date_played is null or (g1.date_played < now() - interval ` +
		`if(WEEKDAY(CURDATE()) >= 5, 7, WEEKDAY(CURDATE()) + 8) ` +
		` day))
		 group by ptg.player_id
		 order by ptg.group_id asc, points desc, setDf desc, u.df desc, p.id asc`
}

func GetTournamentGroupsQuery() string {
	return `select id, name from tournament_group tg where tg.tournament_id = ?`
}

func GetTournamentGroupQuery() string {
	return `select g.play_order as matchNumber, g.id, g.name, 4 as maxStage,
			g.stage, g.home_player_id as hpid, g.away_player_id as apid, g.winner_id,
			g.home_score as homeScoreTotal, g.away_score as awayScoreTotal, g.is_walkover,
			if (g.home_player_id, p1.name, g.playoff_home_player_id) as homePlayerName, 
			if (g.away_player_id, p2.name, g.playoff_away_player_id) as awayPlayerName,
			coalesce(g.level, '') as level, l.name, g.announced, 0 as isFinalGame,
			s1.home_points as s1hp, s1.away_points s1ap,
			s2.home_points as s2hp, s2.away_points s2ap,
			s3.home_points as s3hp, s3.away_points s3ap,
			s4.home_points as s4hp, s4.away_points s4ap,
			s5.home_points as s5hp, s5.away_points s5ap,
			s6.home_points as s5hp, s6.away_points s6ap,
			s7.home_points as s5hp, s7.away_points s7ap
			from game g 
			left join player p1 on p1.id = g.home_player_id 
			left join player p2 on p2.id = g.away_player_id 
			join tournament t on t.id = g.tournament_id and t.is_playoffs = 1
			join tournament_group l on l.id = g.tournament_group_id 
				left join scores s1 on s1.game_id = g.id and s1.set_number = 1
				left join scores s2 on s2.game_id = g.id and s2.set_number = 2
				left join scores s3 on s3.game_id = g.id and s3.set_number = 3
				left join scores s4 on s4.game_id = g.id and s4.set_number = 4
				left join scores s5 on s5.game_id = g.id and s5.set_number = 5
				left join scores s6 on s6.game_id = g.id and s6.set_number = 6
				left join scores s7 on s7.game_id = g.id and s7.set_number = 7
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

func GetTournamentGamesQuery() string {
	return `SELECT
			g.id, tg.name as groupName, g.tournament_id as tournamentId, g.office_id as officeId,
			DATE(g.date_of_match) as dateOfMatch, DATE(g.date_played) as datePlayed,
			coalesce(p1.id, 0) as homePlayerId, coalesce(p2.id, 0) awayPlayerId, 
			coalesce(p1.name, "") homePlayerName, coalesce(p2.name, "") as awayPlayerName,
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

func GetBonusEligibleGamesQuery() string {
return `select tg.name, ph.name, pa.name, ph.slack_name, pa.slack_name,
   g.date_of_match,
   DATEDIFF(CURDATE(), DATE(g.date_of_match)) as days_diff,
   t.enable_timeliness_bonus,
   t.timeliness_bonus_early,
   t.timeliness_bonus_ontime,
   t.timeliness_window_days
from game g
join tournament t on t.id = g.tournament_id
join tournament_group tg on tg.id = g.tournament_group_id
join player ph on ph.id = g.home_player_id
join player pa on pa.id = g.away_player_id
where t.is_playoffs = 0 
  and t.is_finished = 0
  and g.is_finished = 0
  and g.office_id = ?
  and (
      DATE(g.date_of_match) <= CURDATE()
      OR (
          DATE(g.date_of_match) <= CURDATE() + INTERVAL 5 DAY
          AND DAYOFWEEK(g.date_of_match) NOT IN (1, 7)
      )
  )
order by tg.id, g.date_of_match`
}
