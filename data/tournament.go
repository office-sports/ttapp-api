package data

import (
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
	"math"
	"sort"
	"strconv"
	"strings"
)

// GetTournamentPlayersStatistics returns array of player stats
func GetTournamentPlayersStatistics(id int) (*models.TournamentPlayerStatistics, error) {
	t := new(models.TournamentPlayerStatistics)

	var vI, gid int
	var pid int
	var pName, p2Name string
	err := models.DB.QueryRow(queries.GetStatsMostPointsInGameQuery(), id).Scan(&gid, &vI, &pid, &pName)
	if err != nil {
		return nil, err
	}

	t.MostPointsGid = gid
	t.MostPointsInGame = vI
	t.MostPointsInGamePlayerId = pid
	t.MostPointsInGamePlayerName = pName

	err = models.DB.QueryRow(queries.GetStatsLeastPointsLostInGameQuery(), id).Scan(&gid, &vI, &pid, &pName)
	if err != nil {
		return nil, err
	}

	t.LeastPointsGid = gid
	t.LeastPointsLostInGame = vI
	t.LeastPointsLostInGamePlayerId = pid
	t.LeastPointsLostInGamePlayerName = pName

	err = models.DB.QueryRow(queries.GetStatsEloGainQuery(), id).Scan(&gid, &vI, &pid, &pName)
	if err != nil {
		return nil, err
	}

	t.MostEloGainGid = gid
	t.MostEloGain = vI
	t.MostEloGainPlayerId = pid
	t.MostEloGainPlayerName = pName

	err = models.DB.QueryRow(queries.GetStatsEloLostQuery(), id).Scan(&gid, &vI, &pid, &pName)
	if err != nil {
		return nil, err
	}

	t.MostEloLostGid = gid
	t.MostEloLost = vI
	t.MostEloLostPlayerId = pid
	t.MostEloLostPlayerName = pName

	err = models.DB.QueryRow(queries.GetStatsMostPointsGameQuery(), id).Scan(&gid, &vI, &pName, &p2Name)
	if err != nil {
		return nil, err
	}

	t.MostPointsGameGid = gid
	t.MostPointsGame = vI
	t.MostPointsGameHomeName = pName
	t.MostPointsGameAwayName = p2Name

	err = models.DB.QueryRow(queries.GetStatsLeastPointsGameQuery(), id).Scan(&gid, &vI, &pName, &p2Name)
	if err != nil {
		return nil, err
	}

	t.LeastPointsGameGid = gid
	t.LeastPointsGame = vI
	t.LeastPointsGameHomeName = pName
	t.LeastPointsGameAwayName = p2Name

	// Get longest set streak by player
	var hpid, apid, winid, sid, hwon, awon int
	var dom string

	setStreakByPlayerID := make(map[int]int)
	currentSetStreakByPlayerID := make(map[int]int)

	rows, err := models.DB.Query(queries.GetStatsLongestSetStreak(), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rowCount := 0
	maxStreak := 0
	for rows.Next() {
		err := rows.Scan(&gid, &hpid, &apid, &winid, &dom, &sid, &hwon, &awon)
		if err != nil {
			return nil, err
		}

		// If there are no player mappings, add them to both general and current
		if rowCount == 0 {
			setStreakByPlayerID[hpid] = 0
			currentSetStreakByPlayerID[hpid] = 0
		}
		if rowCount == 0 {
			setStreakByPlayerID[apid] = 0
			currentSetStreakByPlayerID[apid] = 0
		}

		if hwon == 1 {
			currentSetStreakByPlayerID[hpid]++
			currentSetStreakByPlayerID[apid] = 0

			if currentSetStreakByPlayerID[hpid] >= setStreakByPlayerID[hpid] {
				setStreakByPlayerID[hpid] = currentSetStreakByPlayerID[hpid]
			}
			if setStreakByPlayerID[hpid] >= maxStreak {
				maxStreak = setStreakByPlayerID[hpid]
			}
		} else if awon == 1 {
			currentSetStreakByPlayerID[apid]++
			currentSetStreakByPlayerID[hpid] = 0

			if currentSetStreakByPlayerID[apid] >= setStreakByPlayerID[apid] {
				setStreakByPlayerID[apid] = currentSetStreakByPlayerID[apid]
			}
			if setStreakByPlayerID[apid] >= maxStreak {
				maxStreak = setStreakByPlayerID[apid]
			}
		}

		rowCount++
	}

	keys := make([]int, 0, len(setStreakByPlayerID))

	for key := range setStreakByPlayerID {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return setStreakByPlayerID[keys[i]] >= setStreakByPlayerID[keys[j]]
	})

	players, err := GetPlayers()
	if err != nil {
		return nil, err
	}

	playersIndexed := make(map[int]*models.Player)
	for _, p := range players {
		playersIndexed[p.ID] = p
	}

	t.MaxSetStreak = maxStreak

	for _, k := range keys {
		if setStreakByPlayerID[k] == maxStreak {
			t.MaxSetStreakPlayers = append(t.MaxSetStreakPlayers, *playersIndexed[k])
		}
	}

	return t, nil
}

// GetTournamentsStatistics returns array of tournament models
func GetTournamentsStatistics() ([]*models.TournamentStatistics, error) {
	rows, err := models.DB.Query(queries.GetTournamentsStatisticsQuery())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tournamentsStatistics := make([]*models.TournamentStatistics, 0)
	for rows.Next() {
		t := new(models.TournamentStatistics)
		err := rows.Scan(&t.Id, &t.Name, &t.Divisions, &t.SetsPlayed, &t.PointsScored, &t.AvgPointsPerMatch)
		if err != nil {
			return nil, err
		}

		tournamentsStatistics = append(tournamentsStatistics, t)
	}

	if err != nil {
		return nil, err
	}

	return tournamentsStatistics, nil
}

// GetTournaments returns array of tournament models
func GetTournaments() ([]*models.Tournament, error) {
	rows, err := models.DB.Query(queries.GetTournamentsQuery())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tournaments := make([]*models.Tournament, 0)
	for rows.Next() {
		t := new(models.Tournament)
		err := rows.Scan(&t.Id, &t.Name, &t.StartTime, &t.IsPlayoffs, &t.OfficeId, &t.Phase, &t.IsFinished,
			&t.ParentTournamentId, &t.Participants, &t.Scheduled, &t.Finished, &t.Sets, &t.Points)
		if err != nil {
			return nil, err
		}

		tournaments = append(tournaments, t)
	}

	if err != nil {
		return nil, err
	}

	return tournaments, nil
}

// GetLiveTournaments returns array of live tournament models
func GetLiveTournaments() ([]*models.Tournament, error) {
	rows, err := models.DB.Query(queries.GetLiveTournamentsQuery())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tournaments := make([]*models.Tournament, 0)
	for rows.Next() {
		t := new(models.Tournament)
		err := rows.Scan(&t.Id, &t.Name, &t.IsFinished, &t.IsPlayoffs,
			&t.StartTime, &t.Phase, &t.OfficeId, &t.ParentTournamentId, &t.Participants, &t.Scheduled, &t.Finished, &t.Sets, &t.Points)
		if err != nil {
			return nil, err
		}

		if t.IsFinished == 0 {
			tournaments = append(tournaments, t)
		}
	}

	if err != nil {
		return nil, err
	}

	return tournaments, nil
}

// GetTournamentById returns tournament model
func GetTournamentById(id int) (*models.Tournament, error) {
	t := new(models.Tournament)
	err := models.DB.QueryRow(queries.GetTournamentByIdQuery(), id).Scan(
		&t.Id, &t.Name, &t.StartTime, &t.IsPlayoffs, &t.OfficeId, &t.Phase, &t.IsFinished,
		&t.ParentTournamentId, &t.Participants, &t.Scheduled, &t.Finished,
		&t.EnableTimelinessBonus, &t.TimelinessBonusEarly, &t.TimelinessBonusOntime, &t.TimelinessWindowHours)

	if err != nil {
		return nil, err
	}

	return t, nil
}

// GetTournamentGroupSchedule returns array of games for requested tournament id
func GetTournamentGroupSchedule(tid int, gid int) ([]*models.Game, error) {
	if gid == 0 {
		gid = 1000
	}
	q := queries.GetBaseTournamentGroupScheduleQuery()

	rows, err := models.DB.Query(q, tid, gid)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.Game, 0)
	for rows.Next() {
		g := new(models.Game)
		err := rows.Scan(&g.TournamentId, &g.OfficeId, &g.MatchId, &g.GroupName, &g.DateOfMatch,
			&g.HomePlayerId, &g.AwayPlayerId, &g.HomePlayerName, &g.AwayPlayerName,
			&g.HomeScoreTotal, &g.AwayScoreTotal, &g.Mode)
		if err != nil {
			return nil, err
		}

		results = append(results, g)
	}

	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetTournamentGroupGames returns array of games for requested tournament id
func GetTournamentGroupGames(oid int) ([]*models.TournamentGroupSchedule, error) {
	rows, err := models.DB.Query(queries.GetTournamentGroupGamesQuery(), oid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Initialize a map to store grouped game schedules
	groupedSchedules := make(map[string][]models.GameSchedule)
	groupOrder := []string{}

	// Loop through the database rows
	for rows.Next() {
		var gn string // Group name
		gs := new(models.GameSchedule)

		// Scan the row data into variables
		err := rows.Scan(&gn, &gs.HomePlayerId, &gs.AwayPlayerId, &gs.HomePlayerName, &gs.AwayPlayerName, &gs.HomePlayerSlackName,
			&gs.AwayPlayerSlackName, &gs.GameWeek, &gs.DateOfMatch, &gs.IsFinished)
		if err != nil {
			return nil, err
		}

		// Append the game schedule to the appropriate group
		if _, exists := groupedSchedules[gn]; !exists {
			groupOrder = append(groupOrder, gn)
		}

		// Append the game schedule to the appropriate group
		groupedSchedules[gn] = append(groupedSchedules[gn], *gs)
	}

	// Convert the groupedSchedules map into a slice of TournamentGroupSchedule
	var tournamentGroups []*models.TournamentGroupSchedule
	for _, groupName := range groupOrder {
		tournamentGroups = append(tournamentGroups, &models.TournamentGroupSchedule{
			Name:         groupName,
			GameSchedule: groupedSchedules[groupName],
		})
	}

	return tournamentGroups, nil
}

// GetTournamentSchedule returns array of games for requested tournament id
func GetTournamentSchedule(tid int, num int) ([]*models.Game, error) {
	if num == 0 {
		num = 1000
	}
	q := queries.GetBaseTournamentScheduleQuery() + ` where g.is_finished = 0 and tt.is_finished = 0 `
	if tid != 0 {
		q = q + ` and g.tournament_id = ? `
	} else {
		q = q + ` and g.tournament_id != ? `
	}

	q += ` group by g.id order by g.date_of_match, g.id asc limit ? `

	rows, err := models.DB.Query(q, tid, num)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.Game, 0)
	for rows.Next() {
		g := new(models.Game)
		err := rows.Scan(&g.TournamentId, &g.OfficeId, &g.MatchId, &g.GroupName, &g.DateOfMatch,
			&g.HomePlayerId, &g.AwayPlayerId, &g.HomePlayerName, &g.AwayPlayerName,
			&g.HomeScoreTotal, &g.AwayScoreTotal, &g.Mode)
		if err != nil {
			return nil, err
		}

		results = append(results, g)
	}

	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetTournamentResults returns array of finished games for requested tournament id
func GetTournamentResults(tid int, num int) ([]*models.GameResult, error) {
	if num == 0 {
		num = 1000
	}

	q := queries.GetTournamentResultsQuery() + ` where g.is_finished = 1 `
	if tid != 0 {
		q = q + ` and g.tournament_id = ? `
	} else {
		q = q + ` and t.is_finished = ? and t.office_id = 1 `
	}

	q += ` group by g.id order by g.date_played DESC limit ? `

	rows, err := models.DB.Query(q, tid, num)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.GameResult, 0)
	for rows.Next() {
		g := new(models.GameResult)
		ss := new(models.GameResultSetScores)
		err := rows.Scan(&g.MatchId, &g.GroupName, &g.TournamentId, &g.OfficeId, &g.DateOfMatch, &g.DatePlayed,
			&g.HomePlayerId, &g.AwayPlayerId, &g.HomePlayerName, &g.AwayPlayerName, &g.WinnerId, &g.HomeScoreTotal,
			&g.AwayScoreTotal, &g.IsWalkover, &g.HomeElo, &g.NewHomeElo, &g.AwayElo, &g.NewAwayElo,
			&g.HomeEloDiff, &g.AwayEloDiff, &g.HasPoints, &ss.S1hp, &ss.S1ap, &ss.S2hp, &ss.S2ap, &ss.S3hp, &ss.S3ap,
			&ss.S4hp, &ss.S4ap, &ss.S5hp, &ss.S5ap, &ss.S6hp, &ss.S6ap, &ss.S7hp, &ss.S7ap)
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

// GetTournamentGames returns array of all games for requested tournament id
func GetTournamentGames(tid int) ([]*models.GameResult, error) {
	q := queries.GetTournamentGamesQuery() +
		` WHERE g.tournament_id = ? group by g.id order by g.id ASC`

	rows, err := models.DB.Query(q, tid)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.GameResult, 0)
	for rows.Next() {
		g := new(models.GameResult)
		ss := new(models.GameResultSetScores)
		err := rows.Scan(&g.MatchId, &g.GroupName, &g.TournamentId, &g.OfficeId, &g.DateOfMatch, &g.DatePlayed,
			&g.HomePlayerId, &g.AwayPlayerId, &g.HomePlayerName, &g.AwayPlayerName, &g.WinnerId, &g.HomeScoreTotal,
			&g.AwayScoreTotal, &g.IsWalkover, &g.HomeElo, &g.NewHomeElo, &g.AwayElo, &g.NewAwayElo,
			&g.HomeEloDiff, &g.AwayEloDiff, &g.HasPoints, &ss.S1hp, &ss.S1ap, &ss.S2hp, &ss.S2ap, &ss.S3hp, &ss.S3ap,
			&ss.S4hp, &ss.S4ap, &ss.S5hp, &ss.S5ap, &ss.S6hp, &ss.S6ap, &ss.S7hp, &ss.S7ap)
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

func GetTournamentStandingsById(id int) (map[int]*models.TournamentGroup, error) {
	// TODO - fetch this value from db, indicated match number of sets per tournament
	var setsPerGame float64 = 3.0

	rows, err := models.DB.Query(queries.GetTournamentStandingsQuery(), id, id, id, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.GroupStandingsPlayer, 0)
	for rows.Next() {
		p := new(models.GroupStandingsPlayer)
		err := rows.Scan(&p.Pos, &p.PosColor, &p.PlayerId, &p.PlayerName, &p.Played, &p.Wins, &p.Draws, &p.Losses,
			&p.Points, &p.SetsFor, &p.SetsAgainst, &p.SetsDiff, &p.RalliesFor, &p.RalliesAgainst, &p.RalliesDiff,
			&p.GroupId, &p.GroupName, &p.GroupAbbreviation, &p.GroupPromotions)
		if err != nil {
			return nil, err
		}

		results = append(results, p)
	}

	groups := make(map[int]*models.TournamentGroup)
	for _, gp := range results {
		tmp := strings.Split(gp.PosColor, ".")
		gid := gp.GroupId
		if groups[gid] == nil {
			gp.Pos = 1
			gp.PosColor = tmp[gp.Pos-1]
			gr := new(models.TournamentGroup)
			gr.Id = gid
			gr.Name = gp.GroupName
			gr.Abbreviation = gp.GroupAbbreviation
			gr.GroupPromotions = gp.GroupPromotions
			gr.Players = append(gr.Players, *gp)
			groups[gid] = gr
		} else {
			gr := groups[gid]
			gp.Pos = len(gr.Players) + 1
			gp.PosColor = tmp[gp.Pos-1]
			gr.Players = append(gr.Players, *gp)
		}
	}

	for i, group := range groups {
		for j, player := range group.Players {
			player.GamesRemaining = len(group.Players) - player.Played - 1

			player.PointsPotentialMin = player.Points +
				float64(player.SetsDiff)/100
			if player.GamesRemaining == 0 {
				player.PointsPotentialMax =
					player.Points + float64(player.GamesRemaining*2) +
						float64(player.SetsDiff)/100
			} else {
				player.PointsPotentialMax =
					player.Points + float64(player.GamesRemaining*2) +
						float64(player.GamesRemaining)*((setsPerGame+float64(player.SetsDiff))/100)
			}

			groups[i].Players[j] = player
		}
	}

	for i, group := range groups {
		for j, player := range group.Players {
			// how many player have < this one's potential min = player's min position
			pc := getGroupsPointsPotentialCount(group, player, false)
			player.PositionMin = len(group.Players) - pc
			pc = getGroupsPointsPotentialCount(group, player, true)
			player.PositionMax = len(group.Players) - pc
			groups[i].Players[j] = player
		}
	}

	if err != nil {
		return nil, err
	}

	return groups, nil
}

func getPlayersTournamentElos(tid int) (map[int]int, map[int]int, map[int][]int, error) {
	rows, err := models.DB.Query(queries.GetPlayersTournamentEloQuery(), tid, tid)

	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	elos := make(map[int]int, 0)
	lelos := make(map[int]int, 0)
	form := make(map[int][]int, 0)

	gid, pid, elo, lelo, winnerId := 0, 0, 0, 0, 0
	var loopElo, loopLelo int

	for rows.Next() {
		err := rows.Scan(&gid, &pid, &elo, &lelo, &winnerId)
		if err != nil {
			return nil, nil, nil, err
		}

		loopElo = elo
		loopLelo = lelo

		if elos[pid] == 0 {
			elos[pid] = loopElo
		}

		lelos[pid] = loopLelo
		form[pid] = append(form[pid], winnerId)
	}

	return elos, lelos, form, nil
}

func GetTournamentPerformanceById(id int) (map[int]*models.TournamentPerformanceGroup, error) {
	rows, err := models.DB.Query(queries.GetTournamentPerformanceQuery(), id, id, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	performances := make([]*models.PlayerPerformance, 0)
	for rows.Next() {
		p := new(models.PlayerPerformance)
		err := rows.Scan(&p.Pos, &p.PlayerId, &p.PlayerName, &p.LastElo, &p.GroupId, &p.GroupName, &p.GroupAbbreviation,
			&p.Won, &p.Draw, &p.Lost, &p.Finished, &p.Unfinished, &p.Performance, &p.Points, &p.TotalPoints)
		if err != nil {
			return nil, err
		}

		performances = append(performances, p)
	}

	groupsEloSum := make(map[int]int)

	pelos, lelos, form, err := getPlayersTournamentElos(id)

	groups := make(map[int]*models.TournamentPerformanceGroup)
	for _, gp := range performances {
		gp.StartingElo = pelos[gp.PlayerId]
		gp.LastElo = lelos[gp.PlayerId]
		gp.Form = form[gp.PlayerId]

		// reverse form so the latest games are first
		for i, j := 0, len(gp.Form)-1; i < j; i, j = i+1, j-1 {
			gp.Form[i], gp.Form[j] = gp.Form[j], gp.Form[i]
		}

		gid := gp.GroupId
		if groups[gid] == nil {
			gp.Pos = 1
			gr := new(models.TournamentPerformanceGroup)
			gr.Id = gid
			gr.Name = gp.GroupName
			gr.Abbreviation = gp.GroupAbbreviation
			gr.Players = append(gr.Players, *gp)
			groups[gid] = gr
			groupsEloSum[gid] += gp.LastElo
		} else {
			gr := groups[gid]
			gp.Pos = len(gr.Players) + 1
			gr.Players = append(gr.Players, *gp)
			groupsEloSum[gid] += gp.LastElo
		}
	}

	for _, g := range groups {
		g.GroupAvgElo = groupsEloSum[g.Id] / len(g.Players)
	}

	if err != nil {
		return nil, err
	}

	return groups, nil
}

func GetPreviousTournamentStandingsById(id int) (map[int]*models.TournamentGroup, error) {
	var setsPerGame float64 = 3.0
	rows, err := models.DB.Query(queries.GetTournamentStandingsDaysQuery(), id, id, id, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.GroupStandingsPlayer, 0)
	for rows.Next() {
		p := new(models.GroupStandingsPlayer)
		err := rows.Scan(&p.Pos, &p.PosColor, &p.PlayerId, &p.PlayerName, &p.Played, &p.Wins, &p.Draws, &p.Losses,
			&p.Points, &p.SetsFor, &p.SetsAgainst, &p.SetsDiff, &p.RalliesFor, &p.RalliesAgainst, &p.RalliesDiff,
			&p.GroupId, &p.GroupName, &p.GroupAbbreviation, &p.GroupPromotions)
		if err != nil {
			return nil, err
		}

		results = append(results, p)
	}

	groups := make(map[int]*models.TournamentGroup)
	for _, gp := range results {
		tmp := strings.Split(gp.PosColor, ".")
		gid := gp.GroupId
		if groups[gid] == nil {
			gp.Pos = 1
			gp.PosColor = tmp[gp.Pos-1]
			gr := new(models.TournamentGroup)
			gr.Id = gid
			gr.Name = gp.GroupName
			gr.Abbreviation = gp.GroupAbbreviation
			gr.GroupPromotions = gp.GroupPromotions
			gr.Players = append(gr.Players, *gp)
			groups[gid] = gr
		} else {
			gr := groups[gid]
			gp.Pos = len(gr.Players) + 1
			gp.PosColor = tmp[gp.Pos-1]
			gr.Players = append(gr.Players, *gp)
		}
	}

	for i, group := range groups {
		for j, player := range group.Players {
			player.GamesRemaining = len(group.Players) - player.Played - 1

			player.PointsPotentialMin = player.Points +
				float64(player.SetsDiff)/100
			if player.GamesRemaining == 0 {
				player.PointsPotentialMax =
					player.Points + float64(player.GamesRemaining*2) +
						float64(player.SetsDiff)/100
			} else {
				player.PointsPotentialMax =
					player.Points + float64(player.GamesRemaining*2) +
						float64(player.GamesRemaining)*(setsPerGame/100)
			}

			groups[i].Players[j] = player
		}
	}

	for i, group := range groups {
		for j, player := range group.Players {
			// how many player have < this one's potential min = player's min position
			pc := getGroupsPointsPotentialCount(group, player, false)
			player.PositionMin = len(group.Players) - pc
			pc = getGroupsPointsPotentialCount(group, player, true)
			player.PositionMax = len(group.Players) - pc
			groups[i].Players[j] = player
		}
	}

	if err != nil {
		return nil, err
	}

	return groups, nil
}

func GetTournamentInfo(id int) ([]*models.GroupInfo, error) {
	currentStandings, _ := GetTournamentStandingsById(id)
	previousStandings, _ := GetPreviousTournamentStandingsById(id)

	recap := make([]*models.GroupInfo, 0)

	for gid, group := range currentStandings {
		info := new(models.GroupInfo)

		info.TopDrop = 0
		info.TopClimb = 0
		info.TopDropPlayerName = ""
		info.TopClimbPlayerName = ""

		cPlayed := 0
		pPlayed := 0
		cRemaining := 0

		for _, player := range group.Players {
			playerInfo := new(models.PlayerInfo)

			previousPlayerData := GetPreviousStandingsPlayerData(previousStandings[gid], player.PlayerId)

			cPlayed += player.Played
			pPlayed += previousPlayerData.Played
			cRemaining += player.GamesRemaining

			playerInfo.Id = player.PlayerId
			playerInfo.Name = player.PlayerName
			playerInfo.PositionCurrent = player.Pos
			playerInfo.PositionPrevious = previousPlayerData.Pos
			playerInfo.PositionMovement = playerInfo.PositionCurrent - playerInfo.PositionPrevious
			playerInfo.PositionMin = player.PositionMin
			playerInfo.PositionMax = player.PositionMax

			if info.PositionCandidates == nil {
				pc := make(map[int]*models.PositionCandidate)
				info.PositionCandidates = pc
			}

			_, positionExists := info.PositionCandidates[player.PositionMax]
			if positionExists {
				candidates := info.PositionCandidates[player.PositionMax]
				candidates.PlayerNames = append(candidates.PlayerNames, player.PlayerName)
				info.PositionCandidates[player.PositionMax] = candidates
			} else {
				//candidatesMap := make(map[int]*models.PositionCandidate, 0)
				candidate := new(models.PositionCandidate)
				candidate.PlayerNames = append(candidate.PlayerNames, player.PlayerName)
				if player.PositionMin == player.PositionMax {
					candidate.Secured = 1
				}

				info.PositionCandidates[player.PositionMax] = candidate
			}

			if playerInfo.PositionMovement > 0 {
				info.PositionsDown += playerInfo.PositionMovement
				if info.TopDrop == 0 {
					info.TopDrop = playerInfo.PositionMovement
					info.TopDropPlayerName = playerInfo.Name
				} else {
					if playerInfo.PositionMovement > info.TopDrop {
						info.TopDrop = playerInfo.PositionMovement
						info.TopDropPlayerName = playerInfo.Name
					} else if playerInfo.PositionMovement == info.TopDrop {
						info.TopDropPlayerName += ", " + playerInfo.Name
					}
				}
			}

			if playerInfo.PositionMovement == 0 {
				info.PositionsStay += 1
			}

			if playerInfo.PositionMovement < 0 {
				move := -1 * playerInfo.PositionMovement
				info.PositionsUp += move
				if info.TopClimb == 0 {
					info.TopClimb = move
					info.TopClimbPlayerName = playerInfo.Name
				} else {
					if move > info.TopClimb {
						info.TopClimb = move
						info.TopClimbPlayerName = playerInfo.Name
					} else if move == info.TopClimb {
						info.TopClimbPlayerName += ", " + playerInfo.Name
					}
				}
			}

			info.PlayerInfo = append(info.PlayerInfo, *playerInfo)
		}

		info.Id = gid
		info.Name = group.Name
		info.GamesPlayed = (cPlayed - pPlayed) / 2
		info.GamesRemaining = cRemaining / 2
		info.GroupPromotions = group.GroupPromotions

		recap = append(recap, info)
	}

	CreateRecapMessage(recap)

	return recap, nil
}

func GetPreviousStandingsPlayerData(tg *models.TournamentGroup, pid int) *models.GroupStandingsPlayer {
	for _, player := range tg.Players {
		if player.PlayerId == pid {
			return &player
		}
	}

	return nil
}

func getGroupsPointsPotentialCount(group *models.TournamentGroup, p models.GroupStandingsPlayer, c bool) int {
	count := 0
	pid := p.PlayerId
	for _, player := range group.Players {
		if player.PlayerId != pid {
			if c == true {
				// we're checking for greater values
				if player.PointsPotentialMin <= p.PointsPotentialMax {
					count++
				}
			} else {
				// we're checking for smaller values
				if player.PointsPotentialMax <= p.PointsPotentialMin {
					count++
				}
			}
		}
	}

	return count
}

// GetTournamentLadders returns array of playoffs groups with matches
func GetTournamentLadders(id int) ([]*models.Ladder, error) {
	rows, err := models.DB.Query(queries.GetTournamentGroupsQuery(), id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]*models.Ladder, 0)
	for rows.Next() {
		l := new(models.Ladder)
		err := rows.Scan(&l.GroupId, &l.GroupName)

		if err != nil {
			return nil, err
		}

		r, err := models.DB.Query(queries.GetTournamentGroupQuery(), id, l.GroupId)
		if err != nil {
			return nil, err
		}

		lastGameOrder := 0
		for r.Next() {
			group := new(models.LadderGroup)
			ss := new(models.GameResultSetScores)
			err = r.Scan(&group.Order, &group.GameId, &group.GameName, &group.MaxStage,
				&group.Stage, &group.HomePlayerId, &group.AwayPlayerId,
				&group.WinnerId, &group.HomeScoreTotal, &group.AwayScoreTotal, &group.IsWalkover,
				&group.HomePlayerName, &group.AwayPlayerName, &group.Level, &group.GroupName, &group.Announced,
				&group.IsFinalGame,
				&ss.S1hp, &ss.S1ap, &ss.S2hp, &ss.S2ap, &ss.S3hp, &ss.S3ap, &ss.S4hp, &ss.S4ap, &ss.S5hp, &ss.S5ap,
				&ss.S6hp, &ss.S6ap, &ss.S7hp, &ss.S7ap)

			SetPlayoffGameScores(group, 1, ss.S1hp, ss.S1ap)
			SetPlayoffGameScores(group, 2, ss.S2hp, ss.S2ap)
			SetPlayoffGameScores(group, 3, ss.S3hp, ss.S3ap)
			SetPlayoffGameScores(group, 4, ss.S4hp, ss.S4ap)
			SetPlayoffGameScores(group, 5, ss.S5hp, ss.S5ap)
			SetPlayoffGameScores(group, 6, ss.S6hp, ss.S6ap)
			SetPlayoffGameScores(group, 7, ss.S7hp, ss.S7ap)

			// set is final game value
			if group.Stage == group.MaxStage && group.Order > lastGameOrder {
				group.IsFinalGame = 1
			}

			if group.HomePlayerId == 0 {
				s := strings.Split(group.HomePlayerName, ".")
				if s[0] == "W" {
					group.HomePlayerName = "Winner, #" + s[1]
				} else if s[0] == "L" {
					group.HomePlayerName = "Loser, #" + s[1]
				} else if s[0] == "G" {
					group.HomePlayerName = "Group, #" + s[1]
				}
			}

			if group.AwayPlayerId == 0 {
				s := strings.Split(group.AwayPlayerName, ".")
				if s[0] == "W" {
					group.AwayPlayerName = "Winner, #" + s[1]
				} else if s[0] == "L" {
					group.AwayPlayerName = "Loser, #" + s[1]
				} else if s[0] == "G" {
					group.AwayPlayerName = "Group, #" + s[1]
				}
			}

			if group.Level != "" {
				group.Level = strings.Replace(group.Level, "|LEAGUE|", group.GroupName, -1)
			}
			if group.WinnerId == 0 {
				group.Level = strings.Replace(group.Level, "|WINNER|", "Winner", -1)
				group.Level = strings.Replace(group.Level, "|LOSER|", "Loser", -1)
			} else {
				if group.WinnerId == group.HomePlayerId {
					group.Level = strings.Replace(group.Level, "|WINNER|", group.HomePlayerName, -1)
					group.Level = strings.Replace(group.Level, "|LOSER|", group.AwayPlayerName, -1)
				}
				if group.WinnerId == group.AwayPlayerId {
					group.Level = strings.Replace(group.Level, "|WINNER|", group.AwayPlayerName, -1)
					group.Level = strings.Replace(group.Level, "|LOSER|", group.HomePlayerName, -1)
				}
			}

			lastGameOrder = group.Order

			l.LadderGroup = append(l.LadderGroup, *group)
		}

		r.Close()

		groups = append(groups, l)
	}

	return groups, nil
}

func SetPlayoffGameScores(g *models.LadderGroup, setNumber int, hs *int, as *int) {
	if hs != nil && as != nil {
		s := new(models.SetScore)
		s.Set = setNumber
		s.Home = *hs
		s.Away = *as
		g.SetScores = append(g.SetScores, *s)
	}
}

func GetTournamentProbabilities(tournamentId int) ([]*models.GameProbabilities, error) {
	rows, err := models.DB.Query(queries.GetTournamentProbabilitiesQuery(), tournamentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make([]*models.GameProbabilities, 0)
	for rows.Next() {
		g := new(models.GameProbabilities)
		var homePlayerId, awayPlayerId int
		var homeElo, awayElo int
		var isFinished, isAbandoned, isStarted, isWalkover int

		err := rows.Scan(&g.Id, &isFinished, &isAbandoned, &isStarted,
			&isWalkover, &homePlayerId, &awayPlayerId, &g.WinnerId, &homeElo, &awayElo)
		if err != nil {
			return nil, err
		}
		g.IsFinished = isFinished == 1
		g.IsAbandoned = isAbandoned == 1
		g.IsStarted = isStarted == 1
		g.IsWalkover = isWalkover == 1
		g.IsPlayersDecided = homePlayerId > 0 && awayPlayerId > 0
		g.Players = []int{homePlayerId, awayPlayerId}

		g.PlayerElo = make(map[string]int)
		g.PlayerElo[strconv.Itoa(homePlayerId)] = homeElo
		g.PlayerElo[strconv.Itoa(awayPlayerId)] = awayElo

		homeProb, awayProb := calculateWinProbability(homeElo, awayElo)

		g.PlayerWinningProbabilities = make(map[string]float64)
		g.PlayerWinningProbabilities[strconv.Itoa(homePlayerId)] = homeProb
		g.PlayerWinningProbabilities[strconv.Itoa(awayPlayerId)] = awayProb

		g.PlayerOdds = make(map[string]float64)
		g.PlayerOdds[strconv.Itoa(homePlayerId)] = calculateOdds(homeProb)
		g.PlayerOdds[strconv.Itoa(awayPlayerId)] = calculateOdds(awayProb)

		games = append(games, g)
	}

	return games, nil
}

func calculateWinProbability(eloA, eloB int) (float64, float64) {
	expectedA := 1.0 / (1.0 + math.Pow(10, float64(eloB-eloA)/800.0))
	expectedB := 1.0 - expectedA
	return math.Round(expectedA*1000) / 1000, math.Round(expectedB*1000) / 1000
}

func calculateOdds(probability float64) float64 {
	if probability == 0 {
		return 0
	}
	return math.Round((1.0/probability)*1000) / 1000
}
