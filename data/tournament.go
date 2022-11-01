package data

import "github.com/office-sports/ttapp-api/models"

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
