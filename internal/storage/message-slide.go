package storage

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
)

func (s lowerThirdsService) createMessageSlide(d *entities.MessageSlide) error {
	s.logger.Debug("createMessageSlide")

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.Exec(
		`INSERT INTO MessageSlides (
		  id, 
		  meeting_id,
		  meeting_role,
		  slide_type,
		  slide_order,
		  primary_text,
		  secondary_text
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		d.MessageSlideID,
		d.MeetingID,
		d.MeetingRole,
		d.SlideType,
		d.SlideOrder,
		d.PrimaryText,
		d.SecondaryText,
	)
	if err != nil {
		s.logger.Error("createMessageSlide Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) deleteMessageSlide(userID string, slideID uuid.UUID) (int64, error) {
	s.logger.Debug("deleteMessageSlide for userID ", userID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE MessageSlides SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ? AND deleted_dt IS NULL`,
		slideID,
	)
	if err != nil {
		s.logger.Error("deleteMessageSlide error ", err)
		return 0, err
	}
	affectedRows, _ := result.RowsAffected()
	return affectedRows, nil
}

func (s lowerThirdsService) getMessageSlideByID(userID string, slideID uuid.UUID) (*entities.MessageSlide, error) {
	s.logger.Debug("getMessageSlideByID for userID ", userID, ", slideID ", slideID)
	var messageSlide entities.MessageSlide
	err := s.MySqlDB.Get(
		&messageSlide,
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
        INNER JOIN MessageSlides s
          ON s.meeting_id = m.id
		  AND s.id = ?
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		slideID,
		userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		s.logger.Error(err)
		return nil, err
	}
	return &messageSlide, nil
}

func (s lowerThirdsService) getMessageSlidesByMeeting(userID string, meetingID uuid.UUID) ([]entities.MessageSlide, error) {
	s.logger.Debug("getMessageSlidesByMeeting for userID ", userID, ", meetingID ", meetingID)
	var messageSlides []entities.MessageSlide
	err := s.MySqlDB.Select(
		&messageSlides,
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
        INNER JOIN MessageSlides s
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
	return messageSlides, nil
}

func (s lowerThirdsService) getMessageSlidesByUser(userID string) ([]entities.MessageSlide, error) {
	s.logger.Debug("getMessageSlidesByUser for userID ", userID)
	var messageSlides []entities.MessageSlide
	err := s.MySqlDB.Select(
		&messageSlides,
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
        INNER JOIN MessageSlides s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return messageSlides, nil
}

func (s lowerThirdsService) updateMessageSlide(messageSlideID uuid.UUID, d *entities.MessageSlide) error {
	s.logger.Debug("updateMessageSlide")

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE MessageSlides SET
		  id = ?,
		  meeting_id = ?,
		  meeting_role = ?,
		  slide_type = ?,
		  slide_order = ?,
		  primary_text = ?,
		  secondary_text = ?
        WHERE id = ?`,
		d.MessageSlideID,
		d.MeetingID,
		d.MeetingRole,
		d.SlideType,
		d.SlideOrder,
		d.PrimaryText,
		d.SecondaryText,
		messageSlideID,
	)
	if err != nil {
		s.logger.Error("updateMessageSlide Error", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("updateMessageSlide affected rows: ", affectedRows)
	}
	return nil
}
