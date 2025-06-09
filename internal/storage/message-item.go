package storage

import (
	"database/sql"
	"errors"
	"lowerthirdsapi/internal/entities"

	"github.com/google/uuid"
)

func (s lowerThirdsService) createMessageItem(d *entities.MessageItem) error {
	s.logger.Debug("createMessageItem")

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.Exec(
		`INSERT INTO MessageItems (
		  id, 
		  meeting_id,
		  meeting_role,
		  item_type,
		  item_order,
		  primary_text,
		  secondary_text
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		d.MessageItemID,
		d.MeetingID,
		d.MeetingRole,
		d.ItemType,
		d.ItemOrder,
		d.PrimaryText,
		d.SecondaryText,
	)
	if err != nil {
		s.logger.Error("createMessageItem Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) deleteMessageItem(userID uuid.UUID, itemID uuid.UUID) (int64, error) {
	s.logger.Debug("deleteMessageItem for userID ", userID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE MessageItems SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ? AND deleted_dt IS NULL`,
		itemID,
	)
	if err != nil {
		s.logger.Error("deleteMessageItem error ", err)
		return 0, err
	}
	affectedRows, _ := result.RowsAffected()
	return affectedRows, nil
}

func (s lowerThirdsService) getMessageItemByID(userID uuid.UUID, itemID uuid.UUID) (*entities.MessageItem, error) {
	s.logger.Debug("getMessageItemByID for userID ", userID, ", itemID ", itemID)
	var messageItem entities.MessageItem
	err := s.MySqlDB.Get(
		&messageItem,
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
        INNER JOIN MessageItems s
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
	return &messageItem, nil
}

func (s lowerThirdsService) getMessageItemsByMeeting(userID uuid.UUID, meetingID uuid.UUID) ([]entities.MessageItem, error) {
	s.logger.Debug("getMessageItemsByMeeting for userID ", userID, ", meetingID ", meetingID)
	var messageItems []entities.MessageItem
	err := s.MySqlDB.Select(
		&messageItems,
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
        INNER JOIN MessageItems s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		meetingID,
		userID)
	if errors.Is(err, sql.ErrNoRows) {
		return []entities.MessageItem{}, nil
	}
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return messageItems, nil
}

func (s lowerThirdsService) getMessageItemsByUser(userID uuid.UUID) ([]entities.MessageItem, error) {
	s.logger.Debug("getMessageItemsByUser for userID ", userID)
	var messageItems []entities.MessageItem
	err := s.MySqlDB.Select(
		&messageItems,
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
        INNER JOIN MessageItems s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return messageItems, nil
}

func (s lowerThirdsService) updateMessageItem(messageItemID uuid.UUID, d *entities.MessageItem) error {
	s.logger.Debug("updateMessageItem")

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE MessageItems SET
		  id = ?,
		  meeting_id = ?,
		  meeting_role = ?,
		  item_type = ?,
		  item_order = ?,
		  primary_text = ?,
		  secondary_text = ?
        WHERE id = ?`,
		d.MessageItemID,
		d.MeetingID,
		d.MeetingRole,
		d.ItemType,
		d.ItemOrder,
		d.PrimaryText,
		d.SecondaryText,
		messageItemID,
	)
	if err != nil {
		s.logger.Error("updateMessageItem Error", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("updateMessageItem Error getting affected rows", err)
		return err
	}
	s.logger.Info("updateMessageItem affected rows: ", affectedRows)
	if affectedRows == 0 {
		return errors.New("sql: no rows in result set")
	}
	return nil
}
