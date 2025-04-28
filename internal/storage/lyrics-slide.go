package storage

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
)

func (s lowerThirdsService) createLyricsSlide(d *entities.LyricsSlide) error {
	s.logger.Debug("createLyricsSlide")

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.Exec(
		`INSERT INTO LyricsSlides (
		  id, 
		  meeting_id,
		  meeting_role,
		  slide_type,
		  slide_order,
		  hymn_id,
		  show_translation
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		d.LyricsSlideID,
		d.MeetingID,
		d.MeetingRole,
		d.SlideType,
		d.SlideOrder,
		d.HymnID,
		d.ShowTranslation,
	)
	if err != nil {
		s.logger.Error("createLyricsSlide Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) deleteLyricsSlide(userID uuid.UUID, slideID uuid.UUID) (int64, error) {
	s.logger.Debug("deleteLyricsSlide for userID ", userID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE LyricsSlides SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ? AND deleted_dt IS NULL`,
		slideID,
	)
	if err != nil {
		s.logger.Error("deleteLyricsSlide error ", err)
		return 0, err
	}
	affectedRows, _ := result.RowsAffected()
	return affectedRows, nil
}

func (s lowerThirdsService) getLyricsSlideByID(userID uuid.UUID, slideID uuid.UUID) (*entities.LyricsSlide, error) {
	s.logger.Debug("getLyricsSlideByID for userID ", userID, ", slideID ", slideID)
	var lyricsSlide entities.LyricsSlide
	err := s.MySqlDB.Get(
		&lyricsSlide,
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
        INNER JOIN LyricsSlides s
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
	return &lyricsSlide, nil
}

func (s lowerThirdsService) getLyricsSlidesByMeeting(userID uuid.UUID, meetingID uuid.UUID) ([]entities.LyricsSlide, error) {
	s.logger.Debug("getLyricsSlidesByMeeting for userID ", userID, ", meetingID ", meetingID)
	var lyricsSlides []entities.LyricsSlide
	err := s.MySqlDB.Select(
		&lyricsSlides,
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
        INNER JOIN LyricsSlides s
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
	return lyricsSlides, nil
}

func (s lowerThirdsService) getLyricsSlidesByUser(userID uuid.UUID) ([]entities.LyricsSlide, error) {
	s.logger.Debug("getLyricsSlidesByUser for userID ", userID)
	var lyricsSlides []entities.LyricsSlide
	err := s.MySqlDB.Select(
		&lyricsSlides,
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
        INNER JOIN LyricsSlides s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return lyricsSlides, nil
}

func (s lowerThirdsService) updateLyricsSlide(lyricsSlideID uuid.UUID, d *entities.LyricsSlide) error {
	s.logger.Debug("updateLyricsSlide")

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE LyricsSlides SET 
		  id = ?,
		  meeting_id = ?,
		  meeting_role = ?,
		  slide_type = ?,
		  slide_order = ?,
		  hymn_id = ?,
		  show_translation = ?
        WHERE id = ?`,
		d.LyricsSlideID,
		d.MeetingID,
		d.MeetingRole,
		d.SlideType,
		d.SlideOrder,
		d.HymnID,
		d.ShowTranslation,
		lyricsSlideID,
	)
	if err != nil {
		s.logger.Error("updateLyricsSlide Error", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("updateLyricsSlide affected rows: ", affectedRows)
	}
	return nil
}
