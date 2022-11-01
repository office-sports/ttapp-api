package data

import (
	"github.com/office-sports/ttapp-api/models"
)

func BasePlayerScoresQuery() string {
	i := `select 
    	g.id, 
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
		p1.tournament_elo as homeElo, 
		p2.tournament_elo as awayElo,
		p1.tournament_elo - p2.tournament_elo as homeEloDiff,
		p2.tournament_elo - p1.tournament_elo as awayEloDiff,
		s1.home_points as s1hp, s1.away_points s1ap,
		s2.home_points as s2hp, s2.away_points s2ap,
		s3.home_points as s3hp, s3.away_points s3ap,
		s4.home_points as s4hp, s4.away_points s4ap,
		s5.home_points as s5hp, s5.away_points s5ap
		from game g
		join game_mode gm on gm.id = g.game_mode_id
		join player p1 on p1.id = g.home_player_id
		join player p2 on p2.id = g.away_player_id
		join tournament_group tg on tg.id = g.tournament_group_id
		left join scores s1 on s1.game_id = g.id and s1.set_number = 1
		left join scores s2 on s2.game_id = g.id and s2.set_number = 2
		left join scores s3 on s3.game_id = g.id and s3.set_number = 3
		left join scores s4 on s4.game_id = g.id and s4.set_number = 4
		left join scores s5 on s5.game_id = g.id and s4.set_number = 5`

	return i
}

func GetPlayers() ([]*models.Player, error) {
	rows, err := models.DB.Query(`
			SELECT p.id, p.name, p.nickname, p.tournament_elo as elo, p.tournament_elo_previous as oldElo, 
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
                order by p.name`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]*models.Player, 0)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.Name, &p.Nickname, &p.Elo, &p.OldElo, &p.EloChange, &p.GamesPlayed,
			&p.Wins, &p.Draws, &p.Losses, &p.OfficeId, &p.WinPercentage)
		if err != nil {
			return nil, err
		}

		players = append(players, p)
	}

	if err != nil {
		return nil, err
	}

	return players, nil
}

func GetPlayerById(id int) (*models.Player, error) {
	p := new(models.Player)
	err := models.DB.QueryRow(`
		select p.id, p.name, count(g.id) as played, p.profile_pic_url as pic, p.tournament_elo as elo, 
		sum(if(p.id = g.winner_id, 1, 0)) as wins, 
		sum(if(g.winner_id = 0, 1, 0)) as draws, 
		sum(if(g.winner_id != 0 and g.winner_id != p.id, 1, 0)) as losses,
		IF (count(g.id) > 0, (sum(if(g.winner_id = p.id, 1, 0)) / count(g.id)) * 100 , 0) as winPercentage
		from player p 
		left join game g on p.id in (g.home_player_id, g.away_player_id) and g.is_finished = 1 
		left join tournament t on g.tournament_id = t.id and t.is_official = 1 
		where p.id = ?  
		group by p.id`, id).
		Scan(&p.ID, &p.Name, &p.GamesPlayed, &p.ProfilePicUrl, &p.Elo, &p.Wins, &p.Draws, &p.Losses, &p.WinPercentage)

	if err != nil {
		return nil, err
	}
	p.NotWinPercentage = 100 - p.WinPercentage
	p.DrawPercentage = float32(p.Draws) / float32(p.GamesPlayed) * 100
	p.NotDrawPercentage = 100 - p.DrawPercentage
	p.LossPercentage = float32(p.Losses) / float32(p.GamesPlayed) * 100
	p.NotLossPercentage = 100 - p.LossPercentage

	rows, err := models.DB.Query(
		`select if (p.id = g.home_player_id, g.new_home_elo, g.new_away_elo) as elo
			from player p
			left join game g on p.id in (g.home_player_id, g.away_player_id)
			join tournament t on g.tournament_id = t.id
			where p.id = ?
			and g.is_finished = 1
			and t.is_official = 1
			order by g.date_played asc, g.id asc`, p.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	e := [2]int{0, 1500}
	p.EloHistory = append(p.EloHistory, e)

	var index = 1
	var elo int
	for rows.Next() {
		err := rows.Scan(&elo)
		if err != nil {
			return nil, err
		}
		e = [2]int{index, elo}

		p.EloHistory = append(p.EloHistory, e)
		index++
	}

	return p, nil
}

func GetPlayerGamesById(pid int, finished int) ([]*models.GameResult, error) {
	rows, err := models.DB.Query(BasePlayerScoresQuery()+
		` where (g.home_player_id = ? OR g.away_player_id = ?)
		and g.is_finished = ? and tg.is_official = 1 
		order by g.date_played desc`, pid, pid, finished)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.GameResult, 0)
	for rows.Next() {
		g := new(models.GameResult)
		ss := new(models.GameResultSetScores)
		err := rows.Scan(&g.MatchId, &g.GroupName, &g.DateOfMatch, &g.DatePlayed, &g.HomePlayerId, &g.AwayPlayerId,
			&g.HomePlayerName, &g.AwayPlayerName, &g.WinnerId, &g.HomeScoreTotal, &g.AwayScoreTotal,
			&g.HomeElo, &g.AwayElo, &g.HomeEloDiff, &g.AwayEloDiff,
			&ss.S1hp, &ss.S1ap, &ss.S2hp, &ss.S2ap, &ss.S3hp, &ss.S3ap, &ss.S4hp, &ss.S4ap, &ss.S5hp, &ss.S5ap)
		if err != nil {
			return nil, err
		}

		results = append(results, g)

		if ss.S1hp != nil {
			s := new(models.SetScore)
			s.Set = 1
			s.Home = *ss.S1hp
			s.Away = *ss.S1ap
			g.SetScores = append(g.SetScores, *s)
		}
		if ss.S2hp != nil {
			s := new(models.SetScore)
			s.Set = 2
			s.Home = *ss.S2hp
			s.Away = *ss.S2ap
			g.SetScores = append(g.SetScores, *s)
		}
		if ss.S3hp != nil {
			s := new(models.SetScore)
			s.Set = 3
			s.Home = *ss.S3hp
			s.Away = *ss.S3ap
			g.SetScores = append(g.SetScores, *s)
		}
		if ss.S4hp != nil {
			s := new(models.SetScore)
			s.Set = 4
			s.Home = *ss.S4hp
			s.Away = *ss.S4ap
			g.SetScores = append(g.SetScores, *s)
		}
		if ss.S5hp != nil {
			s := new(models.SetScore)
			s.Set = 5
			s.Home = *ss.S5hp
			s.Away = *ss.S5ap
			g.SetScores = append(g.SetScores, *s)
		}
	}

	if err != nil {
		return nil, err
	}

	return results, nil
}

//
//select g.id, gm.name, g.winner_id as winnerId, p1.name homePlayerName, p2.name as awayPlayerName,
//p1.tournament_elo as homeElo, p2.tournament_elo as awayElo,
//p1.id as homePlayerId, p2.id awayPlayerId, gm.max_sets as maxSets,
//g.home_score as homeScoreTotal, g.away_score as awayScoreTotal,
//s1.home_points as s1hp, s1.away_points s1ap,
//s2.home_points as s2hp, s2.away_points s2ap,
//s3.home_points as s3hp, s3.away_points s3ap,
//s4.home_points as s4hp, s4.away_points s4ap,
//tg.name as groupName, g.date_of_match as dateOfMatch, g.date_played as datePlayed
//
//public function loadPlayerResults($playerId)
//{
//$matchData = [];
//
//$baseSql = $this->baseQuery();
//$baseSql .= 'where (g.home_player_id = :playerId OR g.away_player_id = :playerId)';
//$baseSql .= 'and g.is_finished = 1 and tg.is_official = 1 ';
//$baseSql .= 'order by g.date_played desc';
//
//$params['playerId'] = $playerId;
//
//$em = $this->getEntityManager();
//$stmt = $em->getConnection()->prepare($baseSql);
//$result = $stmt->executeQuery($params)->fetchAllAssociative();
//
//foreach ($result as $match) {
//$matchId = $match['id'];
//
//$setPoints = [
//$match['s1hp'],
//$match['s1ap'],
//$match['s2hp'],
//$match['s2ap'],
//$match['s3hp'],
//$match['s3ap'],
//$match['s4hp'],
//$match['s4ap'],
//];
//$setPoints = array_filter($setPoints, function ($element) {
//return is_numeric($element);
//});
//$numberOfSets = (int)(count($setPoints) / 2);
//
//$setScores = [];
//for ($i = 1; $i <= $numberOfSets; $i++) {
//$homeScoreVar = 's' . $i . 'hp';
//$awayScoreVar = 's' . $i . 'ap';
//$setScores[] = [
//'set' => $i,
//'home' => $match[$homeScoreVar],
//'away' => $match[$awayScoreVar],
//];
//}
//
//$matchData[] = [
//'matchId' => $matchId,
//'groupName' => $match['groupName'],
//'dateOfMatch' => $match['dateOfMatch'],
//'datePlayed' => $match['datePlayed'],
//'homePlayerId' => $match['homePlayerId'],
//'awayPlayerId' => $match['awayPlayerId'],
//'homePlayerName' => $match['homePlayerName'],
//'awayPlayerName' => $match['awayPlayerName'],
//'winnerId' => $match['winnerId'] ?: 0,
//'homeScoreTotal' => $match['homeScoreTotal'],
//'awayScoreTotal' => $match['awayScoreTotal'],
//'numberOfSets' => $numberOfSets,
//'scores' => $setScores,
//];
//}
//
//return $matchData;
//}
