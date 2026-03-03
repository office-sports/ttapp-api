package data

import (
	"log"

	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
)

func CreateRescheduleRequest(gameID, requesterID int, dateFrom, dateTo string, parentID *int) (int64, error) {
	log.Printf("[RESCHEDULE] CreateRescheduleRequest gameID=%d requesterID=%d dateFrom=%s dateTo=%s", gameID, requesterID, dateFrom, dateTo)
	res, err := models.DB.Exec(queries.InsertRescheduleRequestQuery(), gameID, requesterID, dateFrom, dateTo, gameID)
	if err != nil {
		log.Printf("[RESCHEDULE] CreateRescheduleRequest DB error: %v", err)
		return 0, err
	}
	id, err := res.LastInsertId()
	log.Printf("[RESCHEDULE] CreateRescheduleRequest inserted id=%d err=%v", id, err)
	return id, err
}

func UpdateRescheduleByGameID(gameID, requesterID int, dateFrom, dateTo string) error {
	log.Printf("[RESCHEDULE] UpdateRescheduleByGameID gameID=%d requesterID=%d dateFrom=%s dateTo=%s", gameID, requesterID, dateFrom, dateTo)
	_, err := models.DB.Exec(queries.UpdateRescheduleByGameIDQuery(), requesterID, dateFrom, dateTo, gameID)
	if err != nil {
		log.Printf("[RESCHEDULE] UpdateRescheduleByGameID DB error: %v", err)
	}
	return err
}

func GetRescheduleRequestsForPlayer(playerID int) ([]*models.RescheduleRequest, error) {
	log.Printf("[RESCHEDULE] GetRescheduleRequestsForPlayer playerID=%d", playerID)
	rows, err := models.DB.Query(queries.GetRescheduleRequestsForPlayerQuery(), playerID, playerID)
	if err != nil {
		log.Printf("[RESCHEDULE] GetRescheduleRequestsForPlayer query error: %v", err)
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
			log.Printf("[RESCHEDULE] GetRescheduleRequestsForPlayer scan error: %v", err)
			return nil, err
		}
		result = append(result, r)
	}
	log.Printf("[RESCHEDULE] GetRescheduleRequestsForPlayer returned %d rows", len(result))
	return result, nil
}

func GetRescheduleRequestByID(id int) (*models.RescheduleRequest, error) {
	log.Printf("[RESCHEDULE] GetRescheduleRequestByID id=%d", id)
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
		log.Printf("[RESCHEDULE] GetRescheduleRequestByID scan error: %v", err)
		return nil, err
	}
	return r, nil
}

func GetAllRescheduleRequests() ([]*models.RescheduleRequest, error) {
	log.Printf("[RESCHEDULE] GetAllRescheduleRequests")
	rows, err := models.DB.Query(queries.GetAllRescheduleRequestsQuery())
	if err != nil {
		log.Printf("[RESCHEDULE] GetAllRescheduleRequests query error: %v", err)
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
			log.Printf("[RESCHEDULE] GetAllRescheduleRequests scan error: %v", err)
			return nil, err
		}
		result = append(result, r)
	}
	log.Printf("[RESCHEDULE] GetAllRescheduleRequests returned %d rows", len(result))
	return result, nil
}

func HasPendingRescheduleForGame(gameID int) (bool, error) {
	log.Printf("[RESCHEDULE] HasPendingRescheduleForGame gameID=%d", gameID)
	var count int
	err := models.DB.QueryRow(queries.CountPendingRescheduleForGameQuery(), gameID).Scan(&count)
	if err != nil {
		log.Printf("[RESCHEDULE] HasPendingRescheduleForGame error: %v", err)
	}
	return count > 0, err
}

func UpdateRescheduleStatus(id int, status string) error {
	log.Printf("[RESCHEDULE] UpdateRescheduleStatus id=%d status=%s", id, status)
	_, err := models.DB.Exec(queries.UpdateRescheduleStatusQuery(), status, id)
	if err != nil {
		log.Printf("[RESCHEDULE] UpdateRescheduleStatus DB error: %v", err)
	}
	return err
}

func UpdateGameDateOfMatch(gameID int, date string) error {
	log.Printf("[RESCHEDULE] UpdateGameDateOfMatch gameID=%d date=%s", gameID, date)
	_, err := models.DB.Exec(queries.UpdateGameDateOfMatchQuery(), date, gameID)
	if err != nil {
		log.Printf("[RESCHEDULE] UpdateGameDateOfMatch DB error: %v", err)
	}
	return err
}

