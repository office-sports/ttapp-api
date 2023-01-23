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
				IF (count(g.id) > 0, (sum(if(g.winner_id = p.id, 1, 0)) / count(g.id)) * 100 , 0) as winPercentage
			from player p
			left join game g on p.id in (g.home_player_id, g.away_player_id) and g.is_finished = 1
			left join tournament t on t.id = g.tournament_id and t.is_official = 1
			group by p.id, p.name
			order by p.name`
}

func GetPlayerDataQuery() string {
	return `SELECT 
				p.id, p.name, count(g.id) as played, p.profile_pic_url as pic, p.tournament_elo as elo, 
				sum(if(p.id = g.winner_id, 1, 0)) as wins, 
				sum(if(g.winner_id = 0, 1, 0)) as draws, 
				sum(if(g.winner_id != 0 and g.winner_id != p.id, 1, 0)) as losses,
				IF (count(g.id) > 0, (sum(if(g.winner_id = p.id, 1, 0)) / count(g.id)) * 100 , 0) as winPercentage
			from player p 
			left join game g on p.id in (g.home_player_id, g.away_player_id) and g.is_finished = 1 
			left join tournament t on g.tournament_id = t.id and t.is_official = 1 
			where p.id = ?  
			group by p.id`
}

func GetPlayerEloHistoryQuery() string {
	return `select 
    			if (p.id = g.home_player_id, g.new_home_elo, g.new_away_elo) as elo
			from player p
			left join game g on p.id in (g.home_player_id, g.away_player_id)
			join tournament t on g.tournament_id = t.id
			where p.id = ?
			and g.is_finished = 1
			and t.is_official = 1
			order by g.date_played asc, g.id asc`
}

func GetBasePlayerScoresQuery() string {
	return `select 
    	g.id, 
    	g.tournament_id as tournamentId,
    	g.office_id as officeId,
		tg.name as groupName, 
		g.date_of_match as dateOfMatch, 
		g.date_played as datePlayed,
		p1.id as homePlayerId, 
		p2.id awayPlayerId, 
    	p1.name homePlayerName, 
    	p2.name as awayPlayerName,
    	g.winner_id as winnerId, 
		g.home_score as homeScoreTotal, 
		g.away_score as awayScoreTotal,
		g.is_walkover as isWalkover,
		p1.tournament_elo as homeElo, 
		p2.tournament_elo as awayElo,
		p1.tournament_elo - p2.tournament_elo as homeEloDiff,
		p2.tournament_elo - p1.tournament_elo as awayEloDiff,
		s1.home_points as s1hp, s1.away_points s1ap,
		s2.home_points as s2hp, s2.away_points s2ap,
		s3.home_points as s3hp, s3.away_points s3ap,
		s4.home_points as s4hp, s4.away_points s4ap,
		s5.home_points as s5hp, s5.away_points s5ap,
		s6.home_points as s5hp, s6.away_points s6ap,
		s7.home_points as s5hp, s7.away_points s7ap
		from game g
		join game_mode gm on gm.id = g.game_mode_id
		join player p1 on p1.id = g.home_player_id
		join player p2 on p2.id = g.away_player_id
		join tournament_group tg on tg.id = g.tournament_group_id
		left join scores s1 on s1.game_id = g.id and s1.set_number = 1
		left join scores s2 on s2.game_id = g.id and s2.set_number = 2
		left join scores s3 on s3.game_id = g.id and s3.set_number = 3
		left join scores s4 on s4.game_id = g.id and s4.set_number = 4
		left join scores s5 on s5.game_id = g.id and s5.set_number = 5
		left join scores s6 on s6.game_id = g.id and s6.set_number = 6
		left join scores s7 on s7.game_id = g.id and s7.set_number = 7`
}
