package storage

import (
	"database/sql"
	"errors"
	"lowerthirdsapi/internal/entities"

	"github.com/google/uuid"
)

func (s lowerThirdsService) createTimerItem(d *entities.TimerItem) error {
	s.logger.Debug("createTimerItem")

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.Exec(
		`INSERT INTO TimerItems (
		  id, 
		  meeting_id,
		  meeting_role,
		  item_type,
		  item_order,
		  show_meeting_details
		) VALUES (?, ?, ?, ?, ?, ?)`,
		d.TimerItemID,
		d.MeetingID,
		d.MeetingRole,
		d.ItemType,
		d.ItemOrder,
		d.ShowMeetingDetails,
	)
	if err != nil {
		s.logger.Error("createTimerItem Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) deleteTimerItem(userID uuid.UUID, itemID uuid.UUID) (int64, error) {
	s.logger.Debug("deleteTimerItem for userID ", userID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE TimerItems SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ? AND deleted_dt IS NULL`,
		itemID,
	)
	if err != nil {
		s.logger.Error("deleteTimerItem error ", err)
		return 0, err
	}
	affectedRows, _ := result.RowsAffected()
	return affectedRows, nil
}

func (s lowerThirdsService) getTimerItemByID(userID uuid.UUID, itemID uuid.UUID) (*entities.TimerItem, error) {
	s.logger.Debug("getTimerItemByID for userID ", userID, ", itemID ", itemID)
	var timerItem entities.TimerItem
	err := s.MySqlDB.Get(
		&timerItem,
		`SELECT s.*
        FROM OrgUsers ou
        INNER JOIN Users u
          ON u.id = ou.user_id
          AND u.deleted_dt IS NULL
        INNER JOIN Organization o
          ON o.id = ou.org_id
          AND o.deleted_dt IS NULL
        INNER JOIN Meetings m
          ON m.org_id = ou.org_id
          AND m.deleted_dt IS NULL
        INNER JOIN TimerItems s
          ON s.meeting_id = m.id
		  AND s.id = ?
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		itemID,
		userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		s.logger.Error(err)
		return nil, err
	}
	return &timerItem, nil
}

func (s lowerThirdsService) getTimerItemsByMeeting(userID uuid.UUID, meetingID uuid.UUID) ([]entities.TimerItem, error) {
	s.logger.Debug("getTimerItemsByMeeting for userID ", userID, ", meetingID ", meetingID)
	var timerItems []entities.TimerItem
	err := s.MySqlDB.Select(
		&timerItems,
		`SELECT s.*
        FROM OrgUsers ou
        INNER JOIN Users u
          ON u.id = ou.user_id
          AND u.deleted_dt IS NULL
        INNER JOIN Organization o
          ON o.id = ou.org_id
          AND o.deleted_dt IS NULL
        INNER JOIN Meetings m
          ON m.org_id = ou.org_id
		  AND m.id = ?
          AND m.deleted_dt IS NULL
        INNER JOIN TimerItems s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		meetingID,
		userID)
	if errors.Is(err, sql.ErrNoRows) {
		return []entities.TimerItem{}, nil
	}
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return timerItems, nil
}

func (s lowerThirdsService) getTimerItemsByUser(userID uuid.UUID) ([]entities.TimerItem, error) {
	s.logger.Debug("getTimerItemsByUser for userID ", userID)
	var timerItems []entities.TimerItem
	err := s.MySqlDB.Select(
		&timerItems,
		`SELECT s.*
        FROM OrgUsers ou
        INNER JOIN Users u
          ON u.id = ou.user_id
          AND u.deleted_dt IS NULL
        INNER JOIN Organization o
          ON o.id = ou.org_id
          AND o.deleted_dt IS NULL
        INNER JOIN Meetings m
          ON m.org_id = ou.org_id
          AND m.deleted_dt IS NULL
        INNER JOIN TimerItems s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return timerItems, nil
}

func (s lowerThirdsService) updateTimerItem(timerItemID uuid.UUID, d *entities.TimerItem) error {
	s.logger.Debug("updateTimerItem")

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE TimerItems SET 
		  id = ?,
		  meeting_id = ?,
		  meeting_role = ?,
		  item_type = ?,
		  item_order = ?,
		  show_meeting_details = ?
        WHERE id = ?`,
		d.TimerItemID,
		d.MeetingID,
		d.MeetingRole,
		d.ItemType,
		d.ItemOrder,
		d.ShowMeetingDetails,
		timerItemID,
	)
	if err != nil {
		s.logger.Error("updateTimerItem Error", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("updateTimerItem Error getting affected rows", err)
		return err
	}
	s.logger.Info("updateTimerItem affected rows: ", affectedRows)
	if affectedRows == 0 {
		return errors.New("sql: no rows in result set")
	}
	return nil
}
