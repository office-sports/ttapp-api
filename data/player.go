package data

import (
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
	"log"
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
			&p.Wins, &p.Draws, &p.Losses, &p.OfficeId, &p.WinPercentage, &p.PPS, &p.Active)
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

func GetPlayersAvailability() ([]*models.PlayerAvailability, error) {
	rows, err := models.DB.Query(queries.GetPlayersAvailabilityQuery())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playerMap := make(map[int]*models.PlayerAvailability)
	for rows.Next() {
		var playerId int
		var playerName string
		var gameDate string

		err := rows.Scan(&playerId, &playerName, &gameDate)
		if err != nil {
			return nil, err
		}

		if player, exists := playerMap[playerId]; exists {
			player.GameDates = append(player.GameDates, models.GameDate{Date: gameDate})
		} else {
			playerMap[playerId] = &models.PlayerAvailability{
				ID:        playerId,
				Name:      playerName,
				GameDates: []models.GameDate{{Date: gameDate}},
			}
		}
	}

	players := make([]*models.PlayerAvailability, 0, len(playerMap))
	for _, player := range playerMap {
		players = append(players, player)
	}

	return players, nil
}

// SetPlayerAvailability handles adding player date availability
func SetPlayerAvailability(pd models.PlayerDate) (int64, error) {
	hasDate := 0
	var lid int64 = 0
	err := models.DB.QueryRow(`select count(player_id) from player_availability pa where pa.player_id = ? and pa.date = ?`,
		pd.PlayerId, pd.Date).Scan(&hasDate)

	if err != nil {
		return 0, err
	}

	// allow inserting of new set score only if it does not exist
	if hasDate == 0 && pd.PlayerId != 0 {
		s := `INSERT INTO player_availability (player_id, date) VALUES (?, ?)`
		lid, err := RunTransaction(s, pd.PlayerId, pd.Date)
		return lid, err
	}

	return lid, err
}

func DeletePlayerDate(pd models.PlayerDate) {
	s := `DELETE FROM player_availability WHERE player_id = ? and date = ?`
	args := make([]interface{}, 0)
	args = append(args, pd.PlayerId, pd.Date)
	_, err := RunTransaction(s, args...)

	if err != nil {
		log.Println("Error deleting player date id: ", pd.PlayerId, err)
	}
}

func GetLeaders() ([]*models.Leader, error) {
	rows, err := models.DB.Query(queries.GetLeadersQuery())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leaders := make([]*models.Leader, 0)
	for rows.Next() {
		p := new(models.Leader)

		err := rows.Scan(&p.PlayerId, &p.PlayerName, &p.ProfilePicUrl, &p.OfficeId, &p.GWon, &p.GLost, &p.GDiff,
			&p.PWon, &p.PLost, &p.PDiff, &p.SWon, &p.SLost, &p.SDiff)
		if err != nil {
			return nil, err
		}

		leaders = append(leaders, p)
	}

	if err != nil {
		return nil, err
	}

	return leaders, nil
}

// GetPlayerById returns player model for requested id
func GetPlayerById(id int) (*models.Player, error) {
	p := new(models.Player)
	err := models.DB.QueryRow(queries.GetPlayerDataQuery(), id).
		Scan(&p.ID, &p.Name, &p.GamesPlayed, &p.ProfilePicUrl, &p.Elo, &p.Wins, &p.Draws, &p.Losses, &p.WinPercentage, &p.PPS)

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

// GetPlayerGamesById returns array of player games for requested player id, flag isFinished gets only finished games
func GetPlayerGamesById(pid int, finished int) ([]*models.GameResult, error) {
	rows, err := models.DB.Query(queries.GetGameWithScoresQuery()+
		` where (g.home_player_id = ? OR g.away_player_id = ?)
		and g.is_finished = ? and tg.is_official = 1 and g.is_abandoned = 0
		group by g.id
		order by g.date_played desc`, pid, pid, finished)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.GameResult, 0)
	for rows.Next() {
		g := new(models.GameResult)
		ss := new(models.GameResultSetScores)
		err := rows.Scan(&g.MatchId, &g.MaxSets, &g.WinsRequired, &g.TournamentId, &g.OfficeId, &g.GroupName,
			&g.DateOfMatch, &g.DatePlayed, &g.HomePlayerId, &g.AwayPlayerId, &g.HomePlayerName, &g.AwayPlayerName,
			&g.WinnerId, &g.HomeScoreTotal, &g.AwayScoreTotal, &g.IsWalkover, &g.IsFinished, &g.HomeElo, &g.AwayElo,
			&g.NewHomeElo, &g.NewAwayElo, &g.HomeEloDiff, &g.AwayEloDiff, &ss.S1hp, &ss.S1ap, &ss.S2hp, &ss.S2ap,
			&ss.S3hp, &ss.S3ap, &ss.S4hp, &ss.S4ap, &ss.S5hp, &ss.S5ap, &ss.S6hp, &ss.S6ap, &ss.S7hp, &ss.S7ap,
			&g.CurrentHomePoints, &g.CurrentAwayPoints, &g.CurrentSet, &g.CurrentSetId, &g.HasPoints,
			&g.Announced, &g.TS, &g.Name, &g.PlayOrder, &g.Level)
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

// GetPlayerOpponentsById returns player opponents array
func GetPlayerOpponentsById(id int) ([]*models.PlayerOpponent, error) {
	rows, err := models.DB.Query(queries.GetPlayerOpponentsQuery(), id, id, id, id, id, id, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	opponents := make([]*models.PlayerOpponent, 0)
	for rows.Next() {
		p := new(models.PlayerOpponent)

		err := rows.Scan(&p.Games, &p.OpponentId, &p.OpponentName, &p.Wins, &p.Draws, &p.Losses, &p.PlayerWalkovers, &p.OpponentWalkovers)
		p.Id = id
		if err != nil {
			return nil, err
		}

		opponents = append(opponents, p)
	}

	if err != nil {
		return nil, err
	}

	return opponents, nil
}
