package storage

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
)

func (s lowerThirdsService) createBlankSlide(d *entities.BlankSlide) error {
	s.logger.Debug("createBlankSlide")

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.Exec(
		`INSERT INTO BlankSlides (
		  id, 
		  meeting_id,
		  meeting_role,
		  slide_type,
		  slide_order
		) VALUES (?, ?, ?, ?, ?)`,
		d.BlankSlideID,
		d.MeetingID,
		d.MeetingRole,
		d.SlideType,
		d.SlideOrder,
	)
	if err != nil {
		s.logger.Error("createBlankSlide Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) deleteBlankSlide(userID uuid.UUID, slideID uuid.UUID) (int64, error) {
	s.logger.Debug("deleteBlankSlide for userID ", userID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE BlankSlides SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ? AND deleted_dt IS NULL`,
		slideID,
	)
	if err != nil {
		s.logger.Error("deleteBlankSlide error ", err)
		return 0, err
	}
	affectedRows, _ := result.RowsAffected()
	return affectedRows, nil
}

func (s lowerThirdsService) getBlankSlideByID(userID uuid.UUID, slideID uuid.UUID) (*entities.BlankSlide, error) {
	s.logger.Debug("getBlankSlideByID for userID ", userID, ", slideID ", slideID)
	var blankSlide entities.BlankSlide
	err := s.MySqlDB.Get(
		&blankSlide,
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
        INNER JOIN BlankSlides s
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
	return &blankSlide, nil
}

func (s lowerThirdsService) getBlankSlidesByMeeting(userID uuid.UUID, meetingID uuid.UUID) ([]entities.BlankSlide, error) {
	s.logger.Debug("getBlankSlidesByMeeting for userID ", userID, ", meetingID ", meetingID)
	var blankSlides []entities.BlankSlide
	err := s.MySqlDB.Select(
		&blankSlides,
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
        INNER JOIN BlankSlides s
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
	return blankSlides, nil
}

func (s lowerThirdsService) getBlankSlidesByUser(userID uuid.UUID) ([]entities.BlankSlide, error) {
	s.logger.Debug("getBlankSlidesByUser for userID ", userID)
	var blankSlides []entities.BlankSlide
	err := s.MySqlDB.Select(
		&blankSlides,
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
        INNER JOIN BlankSlides s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return blankSlides, nil
}

func (s lowerThirdsService) updateBlankSlide(blankSlideID uuid.UUID, d *entities.BlankSlide) error {
	s.logger.Debug("updateBlankSlide")

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE BlankSlides SET
		  id = ?,
		  meeting_id = ?,
		  meeting_role = ?,
		  slide_type = ?,
		  slide_order = ?
        WHERE id = ?`,
		d.BlankSlideID,
		d.MeetingID,
		d.MeetingRole,
		d.SlideType,
		d.SlideOrder,
		blankSlideID,
	)
	if err != nil {
		s.logger.Error("updateBlankSlide Error", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("updateBlankSlide affected rows: ", affectedRows)
	}
	return nil
}
