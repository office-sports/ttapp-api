package data

import (
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
)

func CreateRescheduleRequest(gameID, requesterID int, dateFrom, dateTo string, parentID *int) (int64, error) {
	res, err := models.DB.Exec(queries.InsertRescheduleRequestQuery(), gameID, requesterID, dateFrom, dateTo, gameID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func UpdateRescheduleByGameID(gameID, requesterID int, dateFrom, dateTo string) error {
	_, err := models.DB.Exec(queries.UpdateRescheduleByGameIDQuery(), requesterID, dateFrom, dateTo, gameID)
	return err
}

func GetRescheduleRequestsForPlayer(playerID int) ([]*models.RescheduleRequest, error) {
	rows, err := models.DB.Query(queries.GetRescheduleRequestsForPlayerQuery(), playerID, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.RescheduleRequest, 0)
	for rows.Next() {
		r := new(models.RescheduleRequest)
		err := rows.Scan(
			&r.ID, &r.GameID, &r.RequesterID, &r.RequesterName,
			&r.ProposedDateFrom, &r.ProposedDateTo, &r.Status,
			&r.ParentID, &r.CreatedAt, &r.RespondedAt,
			&r.HomePlayerID, &r.AwayPlayerID,
			&r.HomePlayerName, &r.AwayPlayerName,
			&r.CurrentDateOfMatch, &r.OriginalDateOfMatch,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

func GetRescheduleRequestByID(id int) (*models.RescheduleRequest, error) {
	row := models.DB.QueryRow(queries.GetRescheduleRequestByIDQuery(), id)
	r := new(models.RescheduleRequest)
	err := row.Scan(
		&r.ID, &r.GameID, &r.RequesterID, &r.RequesterName,
		&r.ProposedDateFrom, &r.ProposedDateTo, &r.Status,
		&r.ParentID, &r.CreatedAt, &r.RespondedAt,
		&r.HomePlayerID, &r.AwayPlayerID,
		&r.HomePlayerName, &r.AwayPlayerName,
		&r.CurrentDateOfMatch, &r.OriginalDateOfMatch,
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func GetAllRescheduleRequests() ([]*models.RescheduleRequest, error) {
	rows, err := models.DB.Query(queries.GetAllRescheduleRequestsQuery())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*models.RescheduleRequest, 0)
	for rows.Next() {
		r := new(models.RescheduleRequest)
		err := rows.Scan(
			&r.ID, &r.GameID, &r.RequesterID, &r.RequesterName,
			&r.ProposedDateFrom, &r.ProposedDateTo, &r.Status,
			&r.ParentID, &r.CreatedAt, &r.RespondedAt,
			&r.HomePlayerID, &r.AwayPlayerID,
			&r.HomePlayerName, &r.AwayPlayerName,
			&r.CurrentDateOfMatch, &r.OriginalDateOfMatch,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}


func HasPendingRescheduleForGame(gameID int) (bool, error) {
	var count int
	err := models.DB.QueryRow(queries.CountPendingRescheduleForGameQuery(), gameID).Scan(&count)
	return count > 0, err
}


func UpdateRescheduleStatus(id int, status string) error {
	_, err := models.DB.Exec(queries.UpdateRescheduleStatusQuery(), status, id)
	return err
}

func UpdateGameDateOfMatch(gameID int, date string) error {
	_, err := models.DB.Exec(queries.UpdateGameDateOfMatchQuery(), date, gameID)
	return err
}
