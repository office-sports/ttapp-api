package data

import (
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
	"strings"
)

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
			&t.Participants, &t.Scheduled, &t.Finished, &t.Sets, &t.Points)
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
			&t.StartTime, &t.Phase, &t.OfficeId, &t.Participants, &t.Scheduled, &t.Finished, &t.Sets, &t.Points)
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

// GetTournamentById returns tournament model
func GetTournamentById(id int) (*models.Tournament, error) {
	t := new(models.Tournament)
	err := models.DB.QueryRow(queries.GetTournamentByIdQuery(), id).Scan(
		&t.Id, &t.Name, &t.StartTime, &t.IsPlayoffs, &t.OfficeId, &t.Phase, &t.IsFinished,
		&t.Participants, &t.Scheduled, &t.Finished)

	if err != nil {
		return nil, err
	}

	return t, nil
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
			&g.HomeScoreTotal, &g.AwayScoreTotal)
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

func GetTournamentStandingsById(id int) (map[int]*models.TournamentGroup, error) {
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
			&p.GroupId, &p.GroupName, &p.GroupAbbreviation)
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
			gr.Players = append(gr.Players, *gp)
			groups[gid] = gr
		} else {
			gr := groups[gid]
			gp.Pos = len(gr.Players) + 1
			gp.PosColor = tmp[gp.Pos-1]
			gr.Players = append(gr.Players, *gp)
		}
	}

	if err != nil {
		return nil, err
	}

	return groups, nil
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

		for r.Next() {
			group := new(models.LadderGroup)
			ss := new(models.GameResultSetScores)
			err = r.Scan(&group.Order, &group.GameId, &group.GameName, &group.MaxStage,
				&group.Stage, &group.HomePlayerId, &group.AwayPlayerId,
				&group.WinnerId, &group.HomeScoreTotal, &group.AwayScoreTotal, &group.IsWalkover,
				&group.HomePlayerName, &group.AwayPlayerName, &group.Level, &group.GroupName, &group.Announced,
				&ss.S1hp, &ss.S1ap, &ss.S2hp, &ss.S2ap, &ss.S3hp, &ss.S3ap, &ss.S4hp, &ss.S4ap, &ss.S5hp, &ss.S5ap,
				&ss.S6hp, &ss.S6ap, &ss.S7hp, &ss.S7ap)

			SetPlayoffGameScores(group, 1, ss.S1hp, ss.S1ap)
			SetPlayoffGameScores(group, 2, ss.S2hp, ss.S2ap)
			SetPlayoffGameScores(group, 3, ss.S3hp, ss.S3ap)
			SetPlayoffGameScores(group, 4, ss.S4hp, ss.S4ap)
			SetPlayoffGameScores(group, 5, ss.S5hp, ss.S5ap)
			SetPlayoffGameScores(group, 6, ss.S6hp, ss.S6ap)
			SetPlayoffGameScores(group, 7, ss.S7hp, ss.S7ap)

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

			group.Level = strings.Replace(group.Level, "|LEAGUE|", group.GroupName, -1)
			if group.WinnerId == 0 {
				group.Level = strings.Replace(group.Level, "|WINNER|", "Winner", -1)
			} else {
				if group.WinnerId == group.HomePlayerId {
					group.Level = strings.Replace(group.Level, "|WINNER|", group.HomePlayerName, -1)
				}
				if group.WinnerId == group.AwayPlayerId {
					group.Level = strings.Replace(group.Level, "|WINNER|", group.AwayPlayerName, -1)
				}
			}

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
