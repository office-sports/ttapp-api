package queries

func InsertRescheduleRequestQuery() string {
	return `INSERT INTO reschedule_request (game_id, requester_id, proposed_date_from, proposed_date_to, parent_id, original_date_of_match)
			VALUES (?, ?, ?, ?, NULL, (SELECT date_of_match FROM game WHERE id = ?))
			ON DUPLICATE KEY UPDATE
			  requester_id = VALUES(requester_id),
			  proposed_date_from = VALUES(proposed_date_from),
			  proposed_date_to = VALUES(proposed_date_to),
			  original_date_of_match = IF(status = 'accepted', original_date_of_match, (SELECT date_of_match FROM game WHERE id = VALUES(game_id))),
			  status = 'pending',
			  responded_at = NULL,
			  created_at = NOW()`
}

func UpdateRescheduleByGameIDQuery() string {
	return `UPDATE reschedule_request SET requester_id=?, proposed_date_from=?, proposed_date_to=?, status='pending', responded_at=NULL, created_at=NOW() WHERE game_id=?`
}

func GetRescheduleRequestsForPlayerQuery() string {
	return `SELECT
				rr.id, rr.game_id, rr.requester_id, p.name as requester_name,
				rr.proposed_date_from, rr.proposed_date_to, rr.status,
				rr.parent_id, rr.created_at, rr.responded_at,
				g.home_player_id, g.away_player_id,
				ph.name as home_player_name, pa.name as away_player_name,
				COALESCE(g.date_of_match, '') as current_date_of_match,
				COALESCE(rr.original_date_of_match, '') as original_date_of_match
			FROM reschedule_request rr
			JOIN game g ON g.id = rr.game_id
			JOIN player p ON p.id = rr.requester_id
			JOIN player ph ON ph.id = g.home_player_id
			JOIN player pa ON pa.id = g.away_player_id
			WHERE (g.home_player_id = ? OR g.away_player_id = ?)
			ORDER BY rr.created_at DESC`
}

func GetAllRescheduleRequestsQuery() string {
	return `SELECT
				rr.id, rr.game_id, rr.requester_id, p.name as requester_name,
				rr.proposed_date_from, rr.proposed_date_to, rr.status,
				rr.parent_id, rr.created_at, rr.responded_at,
				g.home_player_id, g.away_player_id,
				ph.name as home_player_name, pa.name as away_player_name,
				COALESCE(g.date_of_match, '') as current_date_of_match,
				COALESCE(rr.original_date_of_match, '') as original_date_of_match
			FROM reschedule_request rr
			JOIN game g ON g.id = rr.game_id
			JOIN player p ON p.id = rr.requester_id
			JOIN player ph ON ph.id = g.home_player_id
			JOIN player pa ON pa.id = g.away_player_id
			ORDER BY rr.created_at DESC`
}


func CountPendingRescheduleForGameQuery() string {
	return `SELECT COUNT(*) FROM reschedule_request WHERE game_id = ? AND status = 'pending'`
}


func UpdateRescheduleStatusQuery() string {
	return `UPDATE reschedule_request SET status = ?, responded_at = NOW() WHERE id = ?`
}

func UpdateGameDateOfMatchQuery() string {
	return `UPDATE game SET date_of_match = ? WHERE id = ?`
}

func GetRescheduleRequestByIDQuery() string {
	return `SELECT
				rr.id, rr.game_id, rr.requester_id, p.name as requester_name,
				rr.proposed_date_from, rr.proposed_date_to, rr.status,
				rr.parent_id, rr.created_at, rr.responded_at,
				g.home_player_id, g.away_player_id,
				ph.name as home_player_name, pa.name as away_player_name,
				COALESCE(g.date_of_match, '') as current_date_of_match,
				COALESCE(rr.original_date_of_match, '') as original_date_of_match
			FROM reschedule_request rr
			JOIN game g ON g.id = rr.game_id
			JOIN player p ON p.id = rr.requester_id
			JOIN player ph ON ph.id = g.home_player_id
			JOIN player pa ON pa.id = g.away_player_id
			WHERE rr.id = ?`
}
