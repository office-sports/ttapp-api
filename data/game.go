package data

import (
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
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
		err := rows.Scan(&summary.HomeName, &summary.AwayName, &summary.GroupName, &summary.TournamentName,
			&summary.HomeTotalScore, &summary.AwayTotalScore, &summary.HomeTotalPoints, &summary.AwayTotalPoints,
			&summary.HomePointsPerc, &summary.AwayPointsPerc)
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
		err := rows.Scan(&ge.GameStartingServerId, &ge.IsHomePoint, &ge.IsAwayPoint,
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

		serverId := ge.GameStartingServerId
		var otherServerId int
		if serverId == ge.HomePlayerId {
			otherServerId = ge.AwayPlayerId
			sets[ge.SetNumber].SetSummary.HomeServes++
			if ge.IsHomePoint == 1 {
				summary.HomeOwnServePointsTotal++
			}
		} else {
			otherServerId = ge.HomePlayerId
			sets[ge.SetNumber].SetSummary.AwayServes++
			if ge.IsAwayPoint == 1 {
				summary.AwayOwnServePointsTotal++
			}
		}
		servers := [2]int{serverId, otherServerId}
		pointsScored := homePointsScored + awayPointsScored

		var currentServerIndex int //, setStartingServer, currentserver int

		// check who starts serving in current set
		// 20 points meaning at least one player getting to set ball
		if pointsScored <= 20 {
			currentServerIndex = int(math.Ceil(float64(pointsScored)/2)+float64(ge.SetNumber%2)) % 2
		} else {
			currentServerIndex = int(float64(pointsScored)+float64(ge.SetNumber%2)) % 2
		}
		//setStartingServer := servers[(ge.SetNumber+1)%2]
		currentServer := servers[currentServerIndex]

		if currentServer == ge.HomePlayerId {
			summary.HomeServesTotal++
		} else {
			summary.AwayServesTotal++
		}

	}

	timeline.Summary = *summary
	timeline.Sets = sets

	if err != nil {
		return nil, err
	}

	return timeline, nil
}
