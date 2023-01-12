package data

import (
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
	"strings"
)

func GetTournaments() ([]*models.Tournament, error) {
	rows, err := models.DB.Query(`
			select t.id, t.name, t.start_time, t.is_playoffs, t.office_id,
			IF (t.is_playoffs = 0, 'group', 'playoffs') as phase,
			t.is_finished, count(distinct (g.home_player_id)), count(g.id),
			if(sum(g.is_finished) is null, 0, sum(g.is_finished))
			from tournament t
			left join game g on g.tournament_id = t.id
			where t.is_official = 1
			group by t.id, t.start_time
			order by t.start_time desc`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tournaments := make([]*models.Tournament, 0)
	for rows.Next() {
		t := new(models.Tournament)
		err := rows.Scan(&t.Id, &t.Name, &t.StartTime, &t.IsPlayoffs, &t.OfficeId, &t.Phase, &t.IsFinished,
			&t.Participants, &t.Scheduled, &t.Finished)
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

func GetLiveTournament() ([]*models.Tournament, error) {
	rows, err := models.DB.Query(`
			select t.id as tournamentId, t.name as tournamentName, t.is_finished as isFinished, 
			       t.is_playoffs as isPlayoffs, t.office_id as officeId
                from tournament t
                where t.is_finished = 0 and t.is_official = 1 and t.is_playoffs = 1`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tournaments := make([]*models.Tournament, 0)
	for rows.Next() {
		t := new(models.Tournament)
		err := rows.Scan(&t.Id, &t.Name, &t.IsFinished, &t.IsPlayoffs, &t.OfficeId)
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

func GetTournamentSchedule(tid int, num int) ([]*models.Game, error) {
	q := queries.GetBaseTournamentScheduleQuery() +
		` where g.tournament_id = ? and g.is_finished = 0
		group by g.id order by g.date_of_match, g.id asc limit ?`
	if num == 0 {
		// set max limit of schedules to 1000, thi is way more than needed
		num = 1000
	}

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

func GetTournamentStandingsById(id int) (map[int]*models.TournamentGroup, error) {
	rows, err := models.DB.Query(queries.GetTournamentStandingsQuery(), id, id, id, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*models.GroupStandingsPlayer, 0)
	for rows.Next() {
		p := new(models.GroupStandingsPlayer)
		//ss := new(models.GameResultSetScores)
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
