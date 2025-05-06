package storage

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
)

func (s lowerThirdsService) createTimerSlide(d *entities.TimerSlide) error {
	s.logger.Debug("createTimerSlide")

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.Exec(
		`INSERT INTO TimerSlides (
		  id, 
		  meeting_id,
		  meeting_role,
		  slide_type,
		  slide_order,
		  show_meeting_details
		) VALUES (?, ?, ?, ?, ?, ?)`,
		d.TimerSlideID,
		d.MeetingID,
		d.MeetingRole,
		d.SlideType,
		d.SlideOrder,
		d.ShowMeetingDetails,
	)
	if err != nil {
		s.logger.Error("createTimerSlide Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) deleteTimerSlide(userID string, slideID uuid.UUID) (int64, error) {
	s.logger.Debug("deleteTimerSlide for userID ", userID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE TimerSlides SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ? AND deleted_dt IS NULL`,
		slideID,
	)
	if err != nil {
		s.logger.Error("deleteTimerSlide error ", err)
		return 0, err
	}
	affectedRows, _ := result.RowsAffected()
	return affectedRows, nil
}

func (s lowerThirdsService) getTimerSlideByID(userID string, slideID uuid.UUID) (*entities.TimerSlide, error) {
	s.logger.Debug("getTimerSlideByID for userID ", userID, ", slideID ", slideID)
	var timerSlide entities.TimerSlide
	err := s.MySqlDB.Get(
		&timerSlide,
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
        INNER JOIN TimerSlides s
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
	return &timerSlide, nil
}

func (s lowerThirdsService) getTimerSlidesByMeeting(userID string, meetingID uuid.UUID) ([]entities.TimerSlide, error) {
	s.logger.Debug("getTimerSlidesByMeeting for userID ", userID, ", meetingID ", meetingID)
	var timerSlides []entities.TimerSlide
	err := s.MySqlDB.Select(
		&timerSlides,
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
        INNER JOIN TimerSlides s
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
	return timerSlides, nil
}

func (s lowerThirdsService) getTimerSlidesByUser(userID string) ([]entities.TimerSlide, error) {
	s.logger.Debug("getTimerSlidesByUser for userID ", userID)
	var timerSlides []entities.TimerSlide
	err := s.MySqlDB.Select(
		&timerSlides,
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
        INNER JOIN TimerSlides s
          ON s.meeting_id = m.id
		  AND s.deleted_dt IS NULL
        WHERE ou.user_id = ?
          AND ou.deleted_dt IS NULL`,
		userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return timerSlides, nil
}

func (s lowerThirdsService) updateTimerSlide(timerSlideID uuid.UUID, d *entities.TimerSlide) error {
	s.logger.Debug("updateTimerSlide")

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.Exec(
		`UPDATE TimerSlides SET 
		  id = ?,
		  meeting_id = ?,
		  meeting_role = ?,
		  slide_type = ?,
		  slide_order = ?,
		  show_meeting_details = ?
        WHERE id = ?`,
		d.TimerSlideID,
		d.MeetingID,
		d.MeetingRole,
		d.SlideType,
		d.SlideOrder,
		d.ShowMeetingDetails,
		timerSlideID,
	)
	if err != nil {
		s.logger.Error("updateTimerSlide Error", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("updateTimerSlide affected rows: ", affectedRows)
	}
	return nil
}
