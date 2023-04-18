package data

import (
	"database/sql"
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
	"log"
	"math"
	"strconv"
)

// GetLiveGames returns array of live games models
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

// GetTournamentLiveGames returns array of live games models
func GetTournamentLiveGames(id int) ([]*models.LiveGameData, error) {
	rows, err := models.DB.Query(queries.GetTournamentLiveGamesQuery(), id)

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

// FinalizeGame handles all the game finalizing processes
func FinalizeGame(sf models.SetFinal) {
	// update set scores first
	setScores(sf)
	// set game scores
	increaseGameScore(sf)

	// fetch the game to check if it is finished
	gr, _ := GetGameById(sf.GameId)

	// if any of players reached required number of wins, finish game
	if gr.HomeScoreTotal == sf.WinsRequired || gr.AwayScoreTotal == sf.WinsRequired {
		var winnerId, loserId int
		if gr.HomeScoreTotal > gr.AwayScoreTotal {
			winnerId = gr.HomePlayerId
			loserId = gr.AwayPlayerId
		} else {
			winnerId = gr.AwayPlayerId
			loserId = gr.HomePlayerId
		}

		_, err := RunTransaction(queries.FinishGameQuery(), winnerId, sf.GameId)
		if err != nil {
			log.Println("Error finalizing game, game id: ", sf.GameId, err)
		}

		_, err = RunTransaction(queries.CheckAndFinishTournament(), gr.TournamentId)
		if err != nil {
			log.Println("Error finalizing game, game id: ", sf.GameId, err)
		}

		// if this is playoffs game, update next one, if needed
		if gr.Name != "" && gr.PlayOrder > 0 {
			UpdateNextPlayoffsGame(winnerId, loserId, gr.PlayOrder, gr.TournamentId)
		}

		UpdateGameElo(sf.GameId)
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

func UpdateNextPlayoffsGame(winnerId int, loserId int, playOrder int, tournamentId int) {
	updateStringWinner := "W." + strconv.Itoa(playOrder)
	updateStringLoser := "L." + strconv.Itoa(playOrder)

	_, err := RunTransaction(queries.UpdateNextPlayoffGameHomePlayer(), winnerId, updateStringWinner, tournamentId)
	if err != nil {
		log.Println("Error setting next playoffs game data for player, game id: ", err)
	}

	_, err = RunTransaction(queries.UpdateNextPlayoffGameHomePlayer(), loserId, updateStringLoser, tournamentId)
	if err != nil {
		log.Println("Error setting next playoffs game data for player, game id: ", err)
	}

	_, err = RunTransaction(queries.UpdateNextPlayoffGameAwayPlayer(), winnerId, updateStringWinner, tournamentId)
	if err != nil {
		log.Println("Error setting next playoffs game data for player, game id: ", err)
	}

	_, err = RunTransaction(queries.UpdateNextPlayoffGameAwayPlayer(), loserId, updateStringLoser, tournamentId)
	if err != nil {
		log.Println("Error setting next playoffs game data for player, game id: ", err)
	}
}

// GetGameModes returns array of game modes
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

// GetGameTimeline returns game timeline for requested id
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
		&g.Announced, &g.TS, &g.Name, &g.PlayOrder, &g.Level)

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
		s = `UPDATE game SET home_score = ?, away_score = ?, is_finished = 1, is_walkover = ?,
                date_played = NOW(), winner_id = home_player_id WHERE id = ?`
	} else if as > hs {
		s = `UPDATE game SET home_score = ?, away_score = ?, is_finished = 1, is_walkover = ?,
                date_played = NOW(), winner_id = away_player_id WHERE id = ?`
	} else {
		return
	}
	_, err := RunTransaction(s, hs, as, gr.IsWalkover, gr.GameId)

	if err != nil {
		log.Println("Error updating game data, game id: ", gr.GameId, err)
	}
}

// SaveGameScore handles game closing
func SaveGameScore(gr models.GameSetResults) {
	// delete all points connected to the game when entering final score manually
	deleteGamePoints(gr.GameId)
	// delete previous set scores
	deleteGameScores(gr.GameId)
	// save passed set scores
	saveGameScores(gr)

	hs, as := gr.GetFullScore()

	// save data into game table, set scores, winner, is_finished, date_played
	updateGame(gr, hs, as)

	// update elo
	UpdateGameElo(gr.GameId)

	// fetch the game to check if it is finished
	g, _ := GetGameById(gr.GameId)

	// if this is playoffs game, update next one, if needed
	if g.Name != "" && g.PlayOrder > 0 {
		if g.WinnerId == g.HomePlayerId {
			UpdateNextPlayoffsGame(g.HomePlayerId, g.AwayPlayerId, g.PlayOrder, g.TournamentId)
		} else if g.WinnerId == g.AwayPlayerId {
			UpdateNextPlayoffsGame(g.AwayPlayerId, g.HomePlayerId, g.PlayOrder, g.TournamentId)
		}
	}

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

// AnnounceGame updates db and sends message
func AnnounceGame(gid int) {
	game, err := GetGameById(gid)
	if err != nil {
		log.Println("Error fetching game data for announcement: ", err)
	}

	if game.Announced == 0 {
		SetAnnounced(gid, 1, "0")
		SendStartMessage(game)
	}
}

// UpdateServer sets the server id in the db
func UpdateServer(gr models.ChangeServerPayload) {
	_, err := RunTransaction(queries.UpdateServerQuery(), gr.ServerId, gr.GameId)
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

func IsAnnounced(gid int) (*models.Announcement, error) {
	ann := new(models.Announcement)
	err := models.DB.QueryRow(`select announced, ts from game g where g.id = ?`, gid).Scan(&ann.IsAnnounced, &ann.Ts)

	if err != nil {
		return ann, err
	}

	return ann, nil
}

// RecalculateElo returns array of finished games with ELO values and scores
func RecalculateElo() ([]*models.EloHistory, error) {
	rows, err := models.DB.Query(queries.GetEloHistory())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playersElo := make(map[int]int)
	playersGamesPlayed := make(map[int]int)

	gm := make([]*models.EloHistory, 0)
	for rows.Next() {
		h := new(models.EloHistory)
		err := rows.Scan(&h.GameId, &h.HomePlayerId, &h.AwayPlayerId, &h.HomeScoreTotal, &h.AwayScoreTotal,
			&h.HomeEloOld, &h.AwayEloOld, &h.HomeEloNew, &h.AwayEloNew)
		if err != nil {
			return nil, err
		}

		if playersElo[h.HomePlayerId] == 0 {
			playersElo[h.HomePlayerId] = 1500
			playersGamesPlayed[h.HomePlayerId] = 0
			h.HomeEloOld = 1500
		} else {
			h.HomeEloOld = playersElo[h.HomePlayerId]
			playersGamesPlayed[h.HomePlayerId]++
		}

		h.HomePlayed = playersGamesPlayed[h.HomePlayerId]

		if playersElo[h.AwayPlayerId] == 0 {
			playersElo[h.AwayPlayerId] = 1500
			playersGamesPlayed[h.AwayPlayerId] = 0
			h.AwayEloOld = 1500
		} else {
			h.AwayEloOld = playersElo[h.AwayPlayerId]
			playersGamesPlayed[h.AwayPlayerId]++
		}

		h.AwayPlayed = playersGamesPlayed[h.AwayPlayerId]

		nh, na := CalculateElo(h.HomeEloOld, h.HomePlayed, h.HomeScoreTotal, h.AwayEloOld, h.AwayPlayed, h.AwayScoreTotal)

		h.HomeEloNew = nh
		h.AwayEloNew = na

		playersElo[h.HomePlayerId] = nh
		playersElo[h.AwayPlayerId] = na

		updateGameElos(h.GameId, h.HomeEloOld, h.AwayEloOld, h.HomeEloNew, h.AwayEloNew)
		updatePlayerElo(h.HomePlayerId, h.HomeEloNew, h.HomeEloOld)
		updatePlayerElo(h.AwayPlayerId, h.AwayEloNew, h.AwayEloOld)

		gm = append(gm, h)
	}

	if err != nil {
		return nil, err
	}

	return gm, nil
}

// UpdateGameElo returns array of finished games with ELO values and scores
func UpdateGameElo(id int) {
	gid, err := GetGameById(id)

	if err != nil {
		return
	}

	var ohE, oaE, hp, ap int

	err = models.DB.QueryRow(queries.GetPlayersEloData(), gid.HomePlayerId, gid.HomePlayerId, id).Scan(&hp, &ohE)
	if err == sql.ErrNoRows {
		ohE = 1500
		hp = 0
	} else if err != nil {
		return
	}

	err = models.DB.QueryRow(queries.GetPlayersEloData(), gid.AwayPlayerId, gid.AwayPlayerId, id).Scan(&ap, &oaE)
	if err == sql.ErrNoRows {
		oaE = 1500
		ap = 0
	} else if err != nil {
		return
	}

	nh, na := CalculateElo(ohE, hp, gid.HomeScoreTotal, oaE, ap, gid.AwayScoreTotal)

	updateGameElos(id, ohE, oaE, nh, na)
	updatePlayerElo(gid.HomePlayerId, nh, ohE)
	updatePlayerElo(gid.AwayPlayerId, na, oaE)
}

// CalculateElo will calculate the Elo for each player based on the given information. Returned value is new Elo for player1 and player2 respectively
func CalculateElo(hElo int, hPlayed int, hScore int, aElo int, aPlayed int, aScore int) (int, int) {
	if hPlayed == 0 {
		hPlayed = 1
	}
	if aPlayed == 0 {
		aPlayed = 1
	}

	// P1 = Winner
	// P2 = Looser
	// PD = Points Difference
	// Multiplier = ln(abs(PD) + 1) * (2.2 / ((P1(old)-P2(old)) * 0.001 + 2.2))
	// Elo Winner = P1(old) + 800/num_matches * (1 - 1/(1 + 10 ^ (P2(old) - P1(old) / 400) ) )
	// Elo Looser = P2(old) + 800/num_matches * (0 - 1/(1 + 10 ^ (P2(old) - P1(old) / 400) ) )

	if hScore > aScore {
		multiplier := math.Log(math.Abs(float64(hScore-aScore))+1) * (2.2 / ((float64(hElo-aElo))*0.001 + 2.2))
		hElo, aElo = calculateElo(hElo, hPlayed, aElo, aPlayed, multiplier, false)
	} else if hScore < aScore {
		multiplier := math.Log(math.Abs(float64(hScore-aScore))+1) * (2.2 / ((float64(aElo-hElo))*0.001 + 2.2))
		aElo, hElo = calculateElo(aElo, aPlayed, hElo, hPlayed, multiplier, false)
	} else {
		hElo, aElo = calculateElo(hElo, hPlayed, aElo, aPlayed, 1.0, true)
	}
	// Cap Elo at 400 to avoid players going too low
	if hElo < 400 {
		hElo = 400
	}
	if aElo < 400 {
		aElo = 400
	}
	return hElo, aElo
}

func calculateElo(winnerElo int, winnerGames int, loserElo int, loserGames int, multiplier float64, isDraw bool) (int, int) {
	constant := 800.0

	winner := 1.0
	loser := 0.0
	if isDraw {
		winner = 0.5
		loser = 0.5
	}
	changeWinner := int((constant / float64(winnerGames) * (winner - (1 / (1 + math.Pow(10, float64(loserElo-winnerElo)/400))))) * multiplier)
	calculatedWinner := winnerElo + changeWinner

	changeLooser := int((constant / float64(loserGames) * (loser - (1 / (1 + math.Pow(10, float64(winnerElo-loserElo)/400))))) * multiplier)
	calculatedLooser := loserElo + changeLooser

	return calculatedWinner, calculatedLooser
}

func updateGameElos(g int, ohe int, oae int, nhe int, nae int) {
	var s string
	s = `UPDATE game SET old_home_elo = ?, old_away_elo = ?, new_home_elo = ?, new_away_elo = ? WHERE id = ?`
	_, err := models.DB.Exec(s, ohe, oae, nhe, nae, g)

	if err != nil {
		log.Println("Error updating game ELOS, game id: ", g, err)
	}
}

func updatePlayerElo(pid int, elo int, pelo int) {
	var s string
	s = `UPDATE player p SET p.current_elo = ?, p.tournament_elo = ?, p.tournament_elo_previous = ? WHERE id = ?`
	_, err := models.DB.Exec(s, elo, elo, pelo, pid)

	if err != nil {
		log.Println("Error updating player ELOS, player id: ", pid, err)
	}
}
