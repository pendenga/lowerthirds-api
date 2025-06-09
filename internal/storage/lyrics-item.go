package storage

import (
	"database/sql"
	"errors"
	"lowerthirdsapi/internal/entities"

	"github.com/google/uuid"
)

func (s lowerThirdsService) createLyricsItem(d *entities.LyricsItem) error {
	s.logger.Debug("createLyricsItem")

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.Exec(
		`INSERT INTO LyricsItems (
		  id, 
		  meeting_id,
		  meeting_role,
		  item_type,
		  item_order,
		  hymn_id,
		  show_translation
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		d.LyricsItemID,
		d.MeetingID,
		d.MeetingRole,
		d.ItemType,
		d.ItemOrder,
		d.HymnID,
		d.ShowTranslation,
	)
	if err != nil {
		s.logger.Error("createLyricsItem Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) deleteLyricsItem(userID uuid.UUID, itemID uuid.UUID) (int64, error) {
	s.logger.Debug("deleteLyricsItem for userID ", userID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE LyricsItems SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ? AND deleted_dt IS NULL`,
		itemID,
	)
	if err != nil {
		s.logger.Error("deleteLyricsItem error ", err)
		return 0, err
	}
	affectedRows, _ := result.RowsAffected()
	return affectedRows, nil
}

func (s lowerThirdsService) getLyricsItemByID(userID uuid.UUID, itemID uuid.UUID) (*entities.LyricsItem, error) {
	s.logger.Debug("getLyricsItemByID for userID ", userID, ", itemID ", itemID)
	var lyricsItem entities.LyricsItem
	err := s.MySqlDB.Get(
		&lyricsItem,
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
        INNER JOIN LyricsItems s
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
	return &lyricsItem, nil
}

func (s lowerThirdsService) getLyricsItemsByMeeting(userID uuid.UUID, meetingID uuid.UUID) ([]entities.LyricsItem, error) {
	s.logger.Debug("getLyricsItemsByMeeting for userID ", userID, ", meetingID ", meetingID)
	var lyricsItems []entities.LyricsItem
	err := s.MySqlDB.Select(
		&lyricsItems,
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
        INNER JOIN LyricsItems s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		meetingID,
		userID)
	if errors.Is(err, sql.ErrNoRows) {
		return []entities.LyricsItem{}, nil
	}
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return lyricsItems, nil
}

func (s lowerThirdsService) getLyricsItemsByUser(userID uuid.UUID) ([]entities.LyricsItem, error) {
	s.logger.Debug("getLyricsItemsByUser for userID ", userID)
	var lyricsItems []entities.LyricsItem
	err := s.MySqlDB.Select(
		&lyricsItems,
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
        INNER JOIN LyricsItems s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return lyricsItems, nil
}

func (s lowerThirdsService) updateLyricsItem(lyricsItemID uuid.UUID, d *entities.LyricsItem) error {
	s.logger.Debug("updateLyricsItem")

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE LyricsItems SET 
		  id = ?,
		  meeting_id = ?,
		  meeting_role = ?,
		  item_type = ?,
		  item_order = ?,
		  hymn_id = ?,
		  show_translation = ?
        WHERE id = ?`,
		d.LyricsItemID,
		d.MeetingID,
		d.MeetingRole,
		d.ItemType,
		d.ItemOrder,
		d.HymnID,
		d.ShowTranslation,
		lyricsItemID,
	)
	if err != nil {
		s.logger.Error("updateLyricsItem Error", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("updateLyricsItem Error getting affected rows", err)
		return err
	}
	s.logger.Info("updateLyricsItem affected rows: ", affectedRows)
	if affectedRows == 0 {
		return errors.New("sql: no rows in result set")
	}
	return nil
}
