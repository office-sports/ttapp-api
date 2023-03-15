package data

import (
	"fmt"
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
	"log"
	"math"
)

func GetGameModes() ([]*models.GameMode, error) {
	rows, err := models.DB.Query(queries.GetGameModesQuery())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gm := make([]*models.GameMode, 0)
	for rows.Next() {
		o := new(models.GameMode)
		err := rows.Scan(&o.ID, &o.Name, &o.ShortName, &o.WinsRequired, &o.MaxSets)
		if err != nil {
			return nil, err
		}

		gm = append(gm, o)
	}

	if err != nil {
		return nil, err
	}

	return gm, nil
}

func GetGameTimeline(gid int) (*models.GameTimeline, error) {
	timeline := new(models.GameTimeline)
	summary := new(models.GameTimelineGameSummary)

	// summary
	rows, err := models.DB.Query(queries.GetGameTimelineSummaryQuery()+` where g.id = ? `, gid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&summary.GameStartingServerId, &summary.WinnerId, &summary.HomePlayerId, &summary.AwayPlayerId,
			&summary.HomeName, &summary.AwayName, &summary.GroupName,
			&summary.TournamentName, &summary.HomeTotalScore, &summary.AwayTotalScore, &summary.HomeTotalPoints,
			&summary.AwayTotalPoints, &summary.HomePointsPerc, &summary.AwayPointsPerc)
		if err != nil {
			return nil, err
		}
	}

	// sets
	sets := make(map[int]*models.Set)

	rows, err = models.DB.Query(queries.GetGameTimelineQuery()+` where g.id = ? order by p.id `, gid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	homePointsScored := 0
	awayPointsScored := 0
	for rows.Next() {
		ge := new(models.GameEvent)
		err := rows.Scan(&ge.IsHomePoint, &ge.IsAwayPoint,
			&ge.HomePlayerId, &ge.AwayPlayerId, &ge.Timestamp, &ge.SetNumber)
		if err != nil {
			return nil, err
		}

		set := new(models.Set)
		// if there is no set with current number, create it
		if sets[ge.SetNumber] == nil {
			homePointsScored = 0
			awayPointsScored = 0

			setSummary := new(models.SetSummary)
			set.Events = append(set.Events, ge)
			set.SetSummary = *setSummary

			sets[ge.SetNumber] = set
			set.SetSummary.StartTimestamp = ge.Timestamp
		} else {
			set := sets[ge.SetNumber]
			set.Events = append(set.Events, ge)
		}

		if ge.IsHomePoint == 1 {
			homePointsScored++
			sets[ge.SetNumber].SetSummary.HomePoints++
		} else {
			awayPointsScored++
			sets[ge.SetNumber].SetSummary.AwayPoints++
		}
		ge.HomePointsScored = homePointsScored
		ge.AwayPointsScored = awayPointsScored

		gameFirstServerId := summary.GameStartingServerId
		var gameSecondServerId int
		if gameFirstServerId == ge.HomePlayerId {
			gameSecondServerId = ge.AwayPlayerId
		} else {
			gameSecondServerId = ge.HomePlayerId
		}

		servers := [2]int{gameFirstServerId, gameSecondServerId}
		pointsScored := homePointsScored + awayPointsScored

		var currentServerIndex int //, setStartingServer, currentserver int

		// check who starts serving in current set
		// 20 points meaning at least one player getting to set ball
		if pointsScored <= 20 {
			currentServerIndex = int(math.Ceil(float64(pointsScored)/2)+float64(ge.SetNumber%2)) % 2
		} else {
			currentServerIndex = int(float64(pointsScored)+float64(ge.SetNumber%2)) % 2
		}

		ge.CurrentSetStartingServer = servers[(ge.SetNumber+1)%2]
		ge.CurrentServer = servers[currentServerIndex]

		if ge.CurrentServer == ge.HomePlayerId {
			summary.HomeServesTotal++
			summary.HomeOwnServePointsTotal++
		} else {
			summary.AwayServesTotal++
			summary.AwayOwnServePointsTotal++
		}

		isHomeServer := ge.CurrentServer == ge.HomePlayerId
		if isHomeServer {
			sets[ge.SetNumber].SetSummary.HomeServes++
		}
		if !isHomeServer {
			sets[ge.SetNumber].SetSummary.AwayServes++
		}

		if isHomeServer && ge.IsHomePoint == 1 {
			sets[ge.SetNumber].SetSummary.HomeServePoints++
		}

		if !isHomeServer && ge.IsAwayPoint == 1 {
			sets[ge.SetNumber].SetSummary.AwayServePoints++
		}

		sets[ge.SetNumber].SetSummary.EndTimestamp = ge.Timestamp
	}

	timeline.Summary = *summary
	timeline.Sets = sets

	if err != nil {
		return nil, err
	}

	return timeline, nil
}

func GetGameServe(gid int) (*models.ServeData, error) {
	g := new(models.ServeData)
	err := models.DB.QueryRow(queries.GetGameServeDataQuery(), gid).Scan(
		&g.GameId, &g.SetNumber, &g.FirstGameServer, &g.SecondGameServer,
		&g.CurrentSetFirstServer, &g.CurrentSetSecondServer, &g.CurrentServerId, &g.NumServes)

	if err != nil {
		return nil, err
	}

	return g, nil
}

func GetLiveGames() ([]*models.LiveGameData, error) {
	rows, err := models.DB.Query(queries.GetLiveGamesQuery())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gm := make([]*models.LiveGameData, 0)
	for rows.Next() {
		o := new(models.LiveGameData)
		err := rows.Scan(&o.Id, &o.CurrentSet, &o.HomePlayerName, &o.AwayPlayerName,
			&o.Phase, &o.GroupName)
		if err != nil {
			return nil, err
		}

		gm = append(gm, o)
	}

	if err != nil {
		return nil, err
	}

	return gm, nil

}

func GetGameById(gid int) (*models.GameResult, error) {
	g := new(models.GameResult)
	ss := new(models.GameResultSetScores)
	err := models.DB.QueryRow(queries.GetGameWithScoresQuery()+
		` where g.id = ? order by g.date_played desc`, gid).Scan(
		&g.MatchId, &g.MaxSets, &g.WinsRequired, &g.TournamentId, &g.OfficeId, &g.GroupName, &g.DateOfMatch, &g.DatePlayed,
		&g.HomePlayerId, &g.AwayPlayerId, &g.HomePlayerName, &g.AwayPlayerName, &g.WinnerId, &g.HomeScoreTotal,
		&g.AwayScoreTotal, &g.IsWalkover, &g.IsFinished, &g.HomeElo, &g.AwayElo, &g.NewHomeElo, &g.NewAwayElo, &g.HomeEloDiff, &g.AwayEloDiff,
		&ss.S1hp, &ss.S1ap, &ss.S2hp, &ss.S2ap, &ss.S3hp, &ss.S3ap, &ss.S4hp, &ss.S4ap, &ss.S5hp, &ss.S5ap,
		&ss.S6hp, &ss.S6ap, &ss.S7hp, &ss.S7ap, &g.CurrentHomePoints, &g.CurrentAwayPoints, &g.CurrentSet, &g.CurrentSetId, &g.HasPoints,
		&g.Announced, &g.TS)

	if err != nil {
		return nil, err
	}

	SetGameResultSetScores(g, 1, ss.S1hp, ss.S1ap)
	SetGameResultSetScores(g, 2, ss.S2hp, ss.S2ap)
	SetGameResultSetScores(g, 3, ss.S3hp, ss.S3ap)
	SetGameResultSetScores(g, 4, ss.S4hp, ss.S4ap)
	SetGameResultSetScores(g, 5, ss.S5hp, ss.S5ap)
	SetGameResultSetScores(g, 6, ss.S6hp, ss.S6ap)
	SetGameResultSetScores(g, 7, ss.S7hp, ss.S7ap)

	if err != nil {
		return nil, err
	}

	return g, nil
}

func deleteGamePoints(id int) {
	s := `DELETE FROM points p WHERE p.score_id IN (select s.id from scores s where s.game_id = ?)`
	args := make([]interface{}, 0)
	args = append(args, id)
	_, err := RunTransaction(s, args...)

	if err != nil {
		log.Println("Error deleting game scores, game id: ", id, err)
	}
}

func deleteGameScores(id int) {
	s := `DELETE FROM scores WHERE game_id = ?`
	args := make([]interface{}, 0)
	args = append(args, id)
	_, err := RunTransaction(s, args...)

	if err != nil {
		log.Println("Error deleting game scores, game id: ", id, err)
	}
}

func InsertSetScore(gid int, sn int, hp *int, ap *int) (int64, error) {
	hasSetScore := 0
	var lid int64 = 0
	err := models.DB.QueryRow(`select count(id) from scores s where s.game_id = ? and s.set_number = ?`,
		gid, sn).Scan(&hasSetScore)

	if err != nil {
		return 0, err
	}

	// allow inserting of new set score only if it does not exist
	if hasSetScore == 0 {
		s := `INSERT INTO scores (game_id, set_number, home_points, away_points) VALUES (?, ?, ?, ?)`
		lid, err := RunTransaction(s, gid, sn, hp, ap)
		return lid, err
	}

	return lid, err
}

func saveGameScores(gr models.GameSetResults) {
	if gr.S1hp != nil && gr.S1ap != nil {
		_, _ = InsertSetScore(gr.GameId, 1, gr.S1hp, gr.S1ap)
	}
	if gr.S2hp != nil && gr.S2ap != nil {
		_, _ = InsertSetScore(gr.GameId, 2, gr.S2hp, gr.S2ap)
	}
	if gr.S3hp != nil && gr.S3ap != nil {
		_, _ = InsertSetScore(gr.GameId, 3, gr.S3hp, gr.S3ap)
	}
	if gr.S4hp != nil && gr.S4ap != nil {
		_, _ = InsertSetScore(gr.GameId, 4, gr.S4hp, gr.S4ap)
	}
	if gr.S5hp != nil && gr.S5ap != nil {
		_, _ = InsertSetScore(gr.GameId, 5, gr.S5hp, gr.S5ap)
	}
	if gr.S6hp != nil && gr.S6ap != nil {
		_, _ = InsertSetScore(gr.GameId, 6, gr.S6hp, gr.S6ap)
	}
	if gr.S7hp != nil && gr.S7ap != nil {
		_, _ = InsertSetScore(gr.GameId, 7, gr.S7hp, gr.S7ap)
	}
}

func updateGame(gr models.GameSetResults, hs int, as int) {
	var s string
	if hs > as {
		s = `UPDATE game SET home_score = ?, away_score = ?, is_finished = 1, 
                date_played = NOW(), winner_id = home_player_id WHERE id = ?`
	} else if as > hs {
		s = `UPDATE game SET home_score = ?, away_score = ?, is_finished = 1, 
                date_played = NOW(), winner_id = away_player_id WHERE id = ?`
	} else {
		return
	}
	_, err := RunTransaction(s, hs, as, gr.GameId)

	if err != nil {
		log.Println("Error updating game data, game id: ", gr.GameId, err)
	}
}

func Save(gr models.GameSetResults) {
	deleteGamePoints(gr.GameId)
	deleteGameScores(gr.GameId)
	saveGameScores(gr)
	hs, as := gr.GetFullScore()

	// save data into game table, set scores, winner, is_finished, date_played
	updateGame(gr, hs, as)

	// recalculate elo
	calculateElo(gr.GameId)

	ann, err := IsAnnounced(gr.GameId)
	if err != nil {
		log.Println("Error fetching announcement: ", err)
	}

	if ann.IsAnnounced == 0 {
		game, err := GetGameById(gr.GameId)
		if err != nil {
			log.Println("Error fetching game data for announcement: ", err)
		}
		SendEndSetMessage(game)
		// TODO remove SendFinalMessage(game)
	}
}

func increaseGameScore(sf models.SetFinal) {
	var sql string
	if sf.Home > sf.Away {
		sql = `update game g set home_score = home_score + 1 where g.id = ?`
	} else {
		sql = `update game g set away_score = away_score + 1 where g.id = ?`
	}
	_, err := RunTransaction(sql, sf.GameId)
	if err != nil {
		log.Println("Error updating game score, game id: ", sf.GameId, err)
	}
}

func setScores(sf models.SetFinal) {
	var sql string

	sql = `update scores s set home_points = ?, away_points = ? 
                where s.set_number = (select current_set from game g where g.id = ?)
                and s.game_id = ?`
	_, err := RunTransaction(sql, sf.Home, sf.Away, sf.GameId, sf.GameId)
	if err != nil {
		log.Println("Error updating scores, game id: ", sf.GameId, err)
	}
}

func AnnounceGame(gid int) {
	game, err := GetGameById(gid)
	fmt.Println("loaded game", game.Announced)
	if err != nil {
		log.Println("Error fetching game data for announcement: ", err)
	}

	if game.Announced == 0 {
		SetAnnounced(gid, 1, "0")
		fmt.Println("sending message")
		SendStartMessage(game)
	}
}

func FinalizeGame(sf models.SetFinal) {
	// update set scores first
	setScores(sf)
	// set game scores
	increaseGameScore(sf)

	// fetch the game to check if it is finished
	gr, _ := GetGameById(sf.GameId)

	// if any of players reached required number of wins, finish game
	if gr.HomeScoreTotal == sf.WinsRequired || gr.AwayScoreTotal == sf.WinsRequired {
		var winnerId int
		if gr.HomeScoreTotal > gr.AwayScoreTotal {
			winnerId = gr.HomePlayerId
		} else {
			winnerId = gr.AwayPlayerId
		}

		sql := `update game g set winner_id = ?, current_set = 1, is_finished = 1, date_played = now() where g.id = ?`
		_, err := RunTransaction(sql, winnerId, sf.GameId)
		if err != nil {
			log.Println("Error finalizing game, game id: ", sf.GameId, err)
		}
	} else {
		var z int = 0
		currentSet := gr.CurrentSet + 1
		UpdateGameCurrentSet(sf.GameId, currentSet)
		_, err := InsertSetScore(sf.GameId, currentSet, &z, &z)
		if err != nil {
			log.Println("Error adding set scores, game id: ", sf.GameId, err)
		}
	}

	// reload the game
	gr, _ = GetGameById(sf.GameId)

	SendEndSetMessage(gr)
}

func UpdateServer(gr models.ChangeServerPayload) {
	s := `UPDATE game SET server_id = ? WHERE id = ?`

	_, err := RunTransaction(s, gr.ServerId, gr.GameId)

	if err != nil {
		log.Println("Error updating game server, game id: ", gr.GameId, err)
	}
}

func SetAnnounced(gid int, announced int, ts string) {
	s := `UPDATE game SET announced = ?, ts = ? WHERE id = ?`

	_, err := RunTransaction(s, announced, ts, gid)

	if err != nil {
		log.Println("Error updating game announcement, game id: ", gid, err)
	}
}

func SetTs(gid int, ts string) {
	s := `UPDATE game SET ts = ? WHERE id = ?`

	_, err := RunTransaction(s, ts, gid)

	if err != nil {
		log.Println("Error updating game ts, game id: ", gid, err)
	}
}

func UpdateGameCurrentSet(gid int, sid int) {
	s := `UPDATE game SET current_set = ? WHERE id = ?`

	_, err := RunTransaction(s, sid, gid)

	if err != nil {
		log.Println("Error updating game current set, game id: ", gid, err)
	}
}

func GetEloLastCache() ([]*models.EloCache, error) {
	rows, err := models.DB.Query(queries.GetEloLastCache())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gm := make([]*models.EloCache, 0)
	for rows.Next() {
		o := new(models.EloCache)
		err := rows.Scan(&o.Id, &o.HomePlayerId, &o.AwayPlayerId, &o.WinnerId, &o.HomeScoreTotal, &o.AwayScoreTotal,
			&o.HomeElo, &o.AwayElo, &o.NewHomeElo, &o.NewAwayElo)
		if err != nil {
			return nil, err
		}

		gm = append(gm, o)
	}

	if err != nil {
		return nil, err
	}

	return gm, nil
}

func IsAnnounced(gid int) (*models.Announcement, error) {
	ann := new(models.Announcement)
	err := models.DB.QueryRow(`select announced, ts from game g where g.id = ?`, gid).Scan(&ann.IsAnnounced, &ann.Ts)

	if err != nil {
		return ann, err
	}

	return ann, nil
}

func GetEloCache() ([]*models.EloCache, error) {
	rows, err := models.DB.Query(queries.GetEloCache())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gm := make([]*models.EloCache, 0)
	for rows.Next() {
		o := new(models.EloCache)
		err := rows.Scan(&o.Id, &o.HomePlayerId, &o.AwayPlayerId, &o.WinnerId, &o.HomeScoreTotal, &o.AwayScoreTotal,
			&o.HomeElo, &o.AwayElo, &o.NewHomeElo, &o.NewAwayElo)
		if err != nil {
			return nil, err
		}

		gm = append(gm, o)
	}

	if err != nil {
		return nil, err
	}

	return gm, nil
}

func calculateElo(gameId int) {
	g, err := GetGameById(gameId)
	playerCache := make(map[int]*models.PlayerCache, 0)
	var winner, loser *models.PlayerCache
	homePlayerId := g.HomePlayerId
	awayPlayerId := g.AwayPlayerId
	homeScore := g.HomeScoreTotal
	awayScore := g.AwayScoreTotal
	winnerId := g.WinnerId
	winParam := 1.0
	loseParam := 0.0

	o := new(models.EloCache)
	err = models.DB.QueryRow(queries.GetEloLastCache(), homePlayerId, homePlayerId, homePlayerId, gameId).Scan(
		&o.Id, &o.HomePlayerId, &o.AwayPlayerId, &o.WinnerId, &o.HomeScoreTotal, &o.AwayScoreTotal,
		&o.HomeElo, &o.AwayElo, &o.NewHomeElo, &o.NewAwayElo, &o.GamesPlayed)

	p := new(models.EloCache)
	err = models.DB.QueryRow(queries.GetEloLastCache(), awayPlayerId, awayPlayerId, awayPlayerId, gameId).Scan(
		&p.Id, &p.HomePlayerId, &p.AwayPlayerId, &p.WinnerId, &p.HomeScoreTotal, &p.AwayScoreTotal,
		&p.HomeElo, &p.AwayElo, &p.NewHomeElo, &p.NewAwayElo, &p.GamesPlayed)

	// At this point we have to fetch old home and old away elo
	i := new(models.PlayerCache)
	if o.GamesPlayed == 0 {
		i.Elo = 1500
		i.GamesPlayed = 1
		playerCache[homePlayerId] = i
	} else {
		i.GamesPlayed = o.GamesPlayed
		if homePlayerId == o.HomePlayerId {
			i.Elo = *o.NewHomeElo
		} else {
			i.Elo = *o.NewAwayElo
		}
		playerCache[homePlayerId] = i
	}

	j := new(models.PlayerCache)
	if p.GamesPlayed == 0 {
		j.Elo = 1500
		j.GamesPlayed = 1
		playerCache[awayPlayerId] = j
	} else {
		j.GamesPlayed = p.GamesPlayed
		if awayPlayerId == p.HomePlayerId {
			j.Elo = *p.NewHomeElo
		} else {
			j.Elo = *p.NewAwayElo
		}
		playerCache[awayPlayerId] = j
	}

	oldHomeElo := playerCache[homePlayerId].Elo
	oldAwayElo := playerCache[awayPlayerId].Elo

	if winnerId == homePlayerId {
		winner = playerCache[homePlayerId]
		loser = playerCache[awayPlayerId]
	} else if winnerId == awayPlayerId {
		winner = playerCache[awayPlayerId]
		loser = playerCache[homePlayerId]
	} else {
		winner = playerCache[homePlayerId]
		loser = playerCache[awayPlayerId]
		winParam = 0.5
		loseParam = 0.5
	}

	pointDifference := math.Abs(float64(homeScore - awayScore))
	multiplier := math.Log10(pointDifference+1) * (2.2 / (float64(winner.Elo-loser.Elo)*0.001 + 2.2))

	pow1 := float64(800/winner.GamesPlayed) * (winParam - (1 / (1 + math.Pow10(int(float64(loser.Elo-winner.Elo)/400)))))
	pow2 := float64(800/winner.GamesPlayed) * (loseParam - (1 / (1 + math.Pow10(int(float64(winner.Elo-loser.Elo)/400)))))

	winnerNewElo := float64(winner.Elo) + (pow1 * multiplier)
	loserNewElo := float64(loser.Elo) + (pow2 * multiplier)

	if winnerId == awayPlayerId {
		playerCache[awayPlayerId].Elo = int(winnerNewElo)
		playerCache[awayPlayerId].GamesPlayed++
		playerCache[awayPlayerId].OldElo = oldAwayElo
		playerCache[homePlayerId].Elo = int(loserNewElo)
		playerCache[homePlayerId].GamesPlayed++
		playerCache[homePlayerId].OldElo = oldHomeElo
	} else {
		playerCache[homePlayerId].Elo = int(winnerNewElo)
		playerCache[homePlayerId].GamesPlayed++
		playerCache[homePlayerId].OldElo = oldAwayElo
		playerCache[awayPlayerId].Elo = int(loserNewElo)
		playerCache[awayPlayerId].GamesPlayed++
		playerCache[awayPlayerId].OldElo = oldHomeElo
	}

	updateGameElos(gameId, oldHomeElo, oldAwayElo, playerCache[homePlayerId].Elo, playerCache[awayPlayerId].Elo)

	if err != nil {
		return
	}
}

func recalculateElo() {
	playerCache := make(map[int]*models.PlayerCache, 0)
	var gameId, homePlayerId, awayPlayerId, winnerId, homeScore, awayScore int
	var winner, loser *models.PlayerCache
	//var winParam, loseParam float64

	ec, _ := GetEloCache()

	for _, c := range ec {
		winner = nil
		loser = nil
		gameId = c.Id
		homePlayerId = c.HomePlayerId
		awayPlayerId = c.AwayPlayerId
		winnerId = c.WinnerId
		homeScore = c.HomeScoreTotal
		awayScore = c.AwayScoreTotal
		winParam := 1.0
		loseParam := 0.0

		if playerCache[homePlayerId] == nil {
			i := new(models.PlayerCache)
			i.Elo = 1500
			i.GamesPlayed = 1
			playerCache[homePlayerId] = i
		}

		if playerCache[awayPlayerId] == nil {
			i := new(models.PlayerCache)
			i.Elo = 1500
			i.GamesPlayed = 1
			playerCache[awayPlayerId] = i
		}

		oldHomeElo := playerCache[homePlayerId].Elo
		oldAwayElo := playerCache[awayPlayerId].Elo

		if winnerId == homePlayerId {
			winner = playerCache[homePlayerId]
			loser = playerCache[awayPlayerId]
		} else if winnerId == awayPlayerId {
			winner = playerCache[awayPlayerId]
			loser = playerCache[homePlayerId]
		} else {
			winner = playerCache[homePlayerId]
			loser = playerCache[awayPlayerId]
			winParam = 0.5
			loseParam = 0.5
		}

		pointDifference := math.Abs(float64(homeScore - awayScore))
		multiplier := math.Log10(pointDifference+1) * (2.2 / (float64(winner.Elo-loser.Elo)*0.001 + 2.2))

		pow1 := float64(800/winner.GamesPlayed) * (winParam - (1 / (1 + math.Pow10(int(float64(loser.Elo-winner.Elo)/400)))))
		pow2 := float64(800/winner.GamesPlayed) * (loseParam - (1 / (1 + math.Pow10(int(float64(winner.Elo-loser.Elo)/400)))))

		winnerNewElo := float64(winner.Elo) + (pow1 * multiplier)
		loserNewElo := float64(loser.Elo) + (pow2 * multiplier)

		if winnerId == awayPlayerId {
			playerCache[awayPlayerId].Elo = int(winnerNewElo)
			playerCache[awayPlayerId].GamesPlayed++
			playerCache[awayPlayerId].OldElo = oldAwayElo
			playerCache[homePlayerId].Elo = int(loserNewElo)
			playerCache[homePlayerId].GamesPlayed++
			playerCache[homePlayerId].OldElo = oldHomeElo
		} else {
			playerCache[homePlayerId].Elo = int(winnerNewElo)
			playerCache[homePlayerId].GamesPlayed++
			playerCache[homePlayerId].OldElo = oldAwayElo
			playerCache[awayPlayerId].Elo = int(loserNewElo)
			playerCache[awayPlayerId].GamesPlayed++
			playerCache[awayPlayerId].OldElo = oldHomeElo
		}

		updateGameElos(gameId, oldHomeElo, oldAwayElo, playerCache[homePlayerId].Elo, playerCache[awayPlayerId].Elo)
	}
}

func updateGameElos(g int, ohe int, oae int, nhe int, nae int) {
	var s string
	s = `UPDATE game SET old_home_elo = ?, old_away_elo = ?, new_home_elo = ?, new_away_elo = ? WHERE id = ?`
	_, err := models.DB.Exec(s, ohe, oae, nhe, nae, g)

	if err != nil {
		log.Println("Error updating game ELOS, game id: ", g, err)
	}
}
