package data

import (
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
)

func GetPlayers() ([]*models.Player, error) {
	rows, err := models.DB.Query(queries.GetPlayersDataQuery())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]*models.Player, 0)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.Name, &p.Nickname, &p.Elo, &p.OldElo, &p.EloChange, &p.GamesPlayed,
			&p.Wins, &p.Draws, &p.Losses, &p.OfficeId, &p.WinPercentage, &p.Active)
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
	err := models.DB.QueryRow(queries.GetPlayerDataQuery(), id).
		Scan(&p.ID, &p.Name, &p.GamesPlayed, &p.ProfilePicUrl, &p.Elo, &p.Wins, &p.Draws, &p.Losses, &p.WinPercentage)

	if err != nil {
		return nil, err
	}
	if p.GamesPlayed == 0 {
		p.NotWinPercentage = 0
		p.DrawPercentage = 0
		p.NotDrawPercentage = 0
		p.LossPercentage = 0
		p.NotLossPercentage = 0
	} else {
		p.NotWinPercentage = 100 - p.WinPercentage
		p.DrawPercentage = float32(p.Draws) / float32(p.GamesPlayed) * 100
		p.NotDrawPercentage = 100 - p.DrawPercentage
		p.LossPercentage = float32(p.Losses) / float32(p.GamesPlayed) * 100
		p.NotLossPercentage = 100 - p.LossPercentage
	}

	rows, err := models.DB.Query(queries.GetPlayerEloHistoryQuery(), p.ID)
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
	rows, err := models.DB.Query(queries.GetBasePlayerScoresQuery()+
		` where (g.home_player_id = ? OR g.away_player_id = ?)
		and g.is_finished = ? and tg.is_official = 1 and g.is_abandoned = 0
		order by g.date_played desc`, pid, pid, finished)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.GameResult, 0)
	for rows.Next() {
		g := new(models.GameResult)
		ss := new(models.GameResultSetScores)
		err := rows.Scan(&g.MatchId, &g.MaxSets, &g.TournamentId, &g.OfficeId, &g.GroupName, &g.DateOfMatch, &g.DatePlayed,
			&g.HomePlayerId, &g.AwayPlayerId, &g.HomePlayerName, &g.AwayPlayerName, &g.WinnerId, &g.HomeScoreTotal,
			&g.AwayScoreTotal, &g.IsWalkover, &g.IsFinished, &g.HomeElo, &g.AwayElo, &g.NewHomeElo, &g.NewAwayElo, &g.HomeEloDiff, &g.AwayEloDiff,
			&ss.S1hp, &ss.S1ap, &ss.S2hp, &ss.S2ap, &ss.S3hp, &ss.S3ap, &ss.S4hp, &ss.S4ap, &ss.S5hp, &ss.S5ap,
			&ss.S6hp, &ss.S6ap, &ss.S7hp, &ss.S7ap)
		if err != nil {
			return nil, err
		}

		results = append(results, g)

		SetGameResultSetScores(g, 1, ss.S1hp, ss.S1ap)
		SetGameResultSetScores(g, 2, ss.S2hp, ss.S2ap)
		SetGameResultSetScores(g, 3, ss.S3hp, ss.S3ap)
		SetGameResultSetScores(g, 4, ss.S4hp, ss.S4ap)
		SetGameResultSetScores(g, 5, ss.S5hp, ss.S5ap)
		SetGameResultSetScores(g, 6, ss.S6hp, ss.S6ap)
		SetGameResultSetScores(g, 7, ss.S7hp, ss.S7ap)
	}

	if err != nil {
		return nil, err
	}

	return results, nil
}

func SetGameResultSetScores(g *models.GameResult, setNumber int, hs *int, as *int) {
	if hs != nil && as != nil {
		s := new(models.SetScore)
		s.Set = setNumber
		s.Home = *hs
		s.Away = *as
		g.SetScores = append(g.SetScores, *s)
	}
}
