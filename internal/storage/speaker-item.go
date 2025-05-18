package storage

import (
    "database/sql"
    "errors"
    "github.com/google/uuid"
    "lowerthirdsapi/internal/entities"
)

func (s lowerThirdsService) createSpeakerItem(d *entities.SpeakerItem) error {
    s.logger.Debug("createSpeakerItem")

    // TODO: put some user-level security on this query
    _, err := s.MySqlDB.Exec(
        `INSERT INTO SpeakerItems (
		  id, 
		  meeting_id,
		  meeting_role,
		  item_type,
		  item_order,
		  speaker_name,
		  title
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
        d.SpeakerItemID,
        d.MeetingID,
        d.MeetingRole,
        d.ItemType,
        d.ItemOrder,
        d.SpeakerName,
        d.Title,
    )
    if err != nil {
        s.logger.Error("createSpeakerItem Error", err)
        return err
    }
    return nil
}

func (s lowerThirdsService) deleteSpeakerItem(userID uuid.UUID, itemID uuid.UUID) (int64, error) {
    s.logger.Debug("deleteSpeakerItem for userID ", userID)

    // TODO: put some user level security on this query
    result, err := s.MySqlDB.Exec(
        `UPDATE SpeakerItems SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ? AND deleted_dt IS NULL`,
        itemID,
    )
    if err != nil {
        s.logger.Error("deleteSpeakerItem error ", err)
        return 0, err
    }
    affectedRows, _ := result.RowsAffected()
    return affectedRows, nil
}

func (s lowerThirdsService) getSpeakerItemByID(userID uuid.UUID, itemID uuid.UUID) (*entities.SpeakerItem, error) {
    s.logger.Debug("getSpeakerItemByID for userID ", userID, ", itemID ", itemID)
    var speakerItem entities.SpeakerItem
    err := s.MySqlDB.Get(
        &speakerItem,
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
        INNER JOIN SpeakerItems s
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
    return &speakerItem, nil
}

func (s lowerThirdsService) getSpeakerItemsByMeeting(userID uuid.UUID, meetingID uuid.UUID) ([]entities.SpeakerItem, error) {
    s.logger.Debug("getSpeakerItemsByMeeting for userID ", userID, ", meetingID ", meetingID)
    var speakerItems []entities.SpeakerItem
    err := s.MySqlDB.Select(
        &speakerItems,
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
        INNER JOIN SpeakerItems s
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
    return speakerItems, nil
}

func (s lowerThirdsService) getSpeakerItemsByUser(userID uuid.UUID) ([]entities.SpeakerItem, error) {
    s.logger.Debug("getSpeakerItemsByUser for userID ", userID)
    var speakerItems []entities.SpeakerItem
    err := s.MySqlDB.Select(
        &speakerItems,
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
        INNER JOIN SpeakerItems s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
        userID)
    if err != nil {
        s.logger.Error(err)
        return nil, err
    }
    return speakerItems, nil
}

func (s lowerThirdsService) updateSpeakerItem(speakerItemID uuid.UUID, d *entities.SpeakerItem) error {
    s.logger.Debug("updateSpeakerItem")

    // TODO: put some user level security on this query
    result, err := s.MySqlDB.Exec(
        `UPDATE SpeakerItems SET
		  id = ?,
		  meeting_id = ?,
		  meeting_role = ?,
		  item_type = ?,
		  item_order = ?,
		  speaker_name = ?,
		  title = ?
        WHERE id = ?`,
        d.SpeakerItemID,
        d.MeetingID,
        d.MeetingRole,
        d.ItemType,
        d.ItemOrder,
        d.SpeakerName,
        d.Title,
        speakerItemID,
    )
    if err != nil {
        s.logger.Error("updateSpeakerItem Error", err)
        return err
    }
    affectedRows, err := result.RowsAffected()
    if err == nil {
        s.logger.Info("updateSpeakerItem affected rows: ", affectedRows)
    }
    return nil
}
