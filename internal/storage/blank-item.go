package storage

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
)

func (s lowerThirdsService) createBlankItem(d *entities.BlankItem) error {
	s.logger.Debug("createBlankItem")
	s.logger.Debugf("createBlankItem %+v", d)

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.Exec(
		`INSERT INTO BlankItems (
		  id, 
		  meeting_id,
		  meeting_role,
		  item_type,
		  item_order
		) VALUES (?, ?, ?, ?, ?)`,
		d.BlankItemID,
		d.MeetingID,
		d.MeetingRole,
		d.ItemType,
		d.ItemOrder,
	)
	if err != nil {
		s.logger.Error("[createBlankItem] Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) deleteBlankItem(userID uuid.UUID, itemID uuid.UUID) (int64, error) {
	s.logger.Debug("deleteBlankItem for userID ", userID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE BlankItems SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ? AND deleted_dt IS NULL`,
		itemID,
	)
	if err != nil {
		s.logger.Error("[deleteBlankItem] error ", err)
		return 0, err
	}
	affectedRows, _ := result.RowsAffected()
	return affectedRows, nil
}

func (s lowerThirdsService) getBlankItemByID(userID uuid.UUID, itemID uuid.UUID) (*entities.BlankItem, error) {
	s.logger.Debug("getBlankItemByID for userID ", userID, ", itemID ", itemID)
	var blankItem entities.BlankItem
	err := s.MySqlDB.Get(
		&blankItem,
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
        INNER JOIN BlankItems s
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
	return &blankItem, nil
}

func (s lowerThirdsService) getBlankItemsByMeeting(userID uuid.UUID, meetingID uuid.UUID) ([]entities.BlankItem, error) {
	s.logger.Debug("getBlankItemsByMeeting for userID ", userID, ", meetingID ", meetingID)
	var blankItems []entities.BlankItem
	err := s.MySqlDB.Select(
		&blankItems,
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
        INNER JOIN BlankItems s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		meetingID,
		userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return blankItems, nil
}

func (s lowerThirdsService) getBlankItemsByUser(userID uuid.UUID) ([]entities.BlankItem, error) {
	s.logger.Debug("getBlankItemsByUser for userID ", userID)
	var blankItems []entities.BlankItem
	err := s.MySqlDB.Select(
		&blankItems,
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
        INNER JOIN BlankItems s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return blankItems, nil
}

func (s lowerThirdsService) updateBlankItem(blankItemID uuid.UUID, d *entities.BlankItem) error {
	s.logger.Debug("[updateBlankItem]")

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE BlankItems SET
		  id = ?,
		  meeting_id = ?,
		  meeting_role = ?,
		  item_type = ?,
		  item_order = ?
        WHERE id = ?`,
		d.BlankItemID,
		d.MeetingID,
		d.MeetingRole,
		d.ItemType,
		d.ItemOrder,
		blankItemID,
	)
	if err != nil {
		s.logger.Error("[updateBlankItem] Error", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("[updateBlankItem] affected rows: ", affectedRows)
	}
	return nil
}
