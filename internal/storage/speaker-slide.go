package storage

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
)

func (s lowerThirdsService) createSpeakerSlide(d *entities.SpeakerSlide) error {
	s.logger.Debug("createSpeakerSlide")

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.Exec(
		`INSERT INTO SpeakerSlides (
		  id, 
		  meeting_id,
		  meeting_role,
		  slide_type,
		  slide_order,
		  speaker_name,
		  title
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		d.SpeakerSlideID,
		d.MeetingID,
		d.MeetingRole,
		d.SlideType,
		d.SlideOrder,
		d.SpeakerName,
		d.Title,
	)
	if err != nil {
		s.logger.Error("createSpeakerSlide Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) deleteSpeakerSlide(userID string, slideID uuid.UUID) (int64, error) {
	s.logger.Debug("deleteSpeakerSlide for userID ", userID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE SpeakerSlides SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ? AND deleted_dt IS NULL`,
		slideID,
	)
	if err != nil {
		s.logger.Error("deleteSpeakerSlide error ", err)
		return 0, err
	}
	affectedRows, _ := result.RowsAffected()
	return affectedRows, nil
}

func (s lowerThirdsService) getSpeakerSlideByID(userID string, slideID uuid.UUID) (*entities.SpeakerSlide, error) {
	s.logger.Debug("getSpeakerSlideByID for userID ", userID, ", slideID ", slideID)
	var speakerSlide entities.SpeakerSlide
	err := s.MySqlDB.Get(
		&speakerSlide,
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
        INNER JOIN SpeakerSlides s
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
	return &speakerSlide, nil
}

func (s lowerThirdsService) getSpeakerSlidesByMeeting(userID string, meetingID uuid.UUID) ([]entities.SpeakerSlide, error) {
	s.logger.Debug("getSpeakerSlidesByMeeting for userID ", userID, ", meetingID ", meetingID)
	var speakerSlides []entities.SpeakerSlide
	err := s.MySqlDB.Select(
		&speakerSlides,
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
        INNER JOIN SpeakerSlides s
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
	return speakerSlides, nil
}

func (s lowerThirdsService) getSpeakerSlidesByUser(userID string) ([]entities.SpeakerSlide, error) {
	s.logger.Debug("getSpeakerSlidesByUser for userID ", userID)
	var speakerSlides []entities.SpeakerSlide
	err := s.MySqlDB.Select(
		&speakerSlides,
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
        INNER JOIN SpeakerSlides s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return speakerSlides, nil
}

func (s lowerThirdsService) updateSpeakerSlide(speakerSlideID uuid.UUID, d *entities.SpeakerSlide) error {
	s.logger.Debug("updateSpeakerSlide")

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE SpeakerSlides SET
		  id = ?,
		  meeting_id = ?,
		  meeting_role = ?,
		  slide_type = ?,
		  slide_order = ?,
		  speaker_name = ?,
		  title = ?
        WHERE id = ?`,
		d.SpeakerSlideID,
		d.MeetingID,
		d.MeetingRole,
		d.SlideType,
		d.SlideOrder,
		d.SpeakerName,
		d.Title,
		speakerSlideID,
	)
	if err != nil {
		s.logger.Error("updateSpeakerSlide Error", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("updateSpeakerSlide affected rows: ", affectedRows)
	}
	return nil
}
