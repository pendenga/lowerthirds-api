package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
)

func (s lowerThirdsService) CreateMeeting(ctx context.Context, m *entities.Meeting) error {
	s.logger.Debug("CreateMeeting")

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.ExecContext(
		ctx,
		`INSERT INTO Meetings (
			id, org_id, conference, meeting, meeting_date, duration
		) VALUES (?, ?, ?, ?, ?, ?)`,
		m.MeetingID,
		m.OrgID,
		m.Conference,
		m.Meeting,
		m.MeetingDate,
		m.Duration,
	)
	if err != nil {
		s.logger.Error("CreateMeeting Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) DeleteMeeting(ctx context.Context, meetingID uuid.UUID) error {
	s.logger.Debug("DeleteMeeting for meetingID ", meetingID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.ExecContext(ctx, `
		UPDATE Meetings 
		SET deleted_dt = CURRENT_TIMESTAMP 
		WHERE id = ?
		  AND deleted_dt IS NULL`,
		meetingID,
	)
	if err != nil {
		s.logger.Error("DeleteMeeting error ", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("DeleteMeeting affected rows: ", affectedRows)
	}
	return nil
}

func (s lowerThirdsService) GetMeeting(ctx context.Context, meetingID uuid.UUID) (*entities.Meeting, error) {
	userIDVal := ctx.Value(helpers.UserIDKey)
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("userID not found or invalid in context")
	}
	s.logger.Debug("GetMeeting for userID ", userID, " meetingID ", meetingID)
	var meeting entities.Meeting
	err := s.MySqlDB.Get(
		&meeting,
		`SELECT m.*
		FROM OrgUsers ou
		INNER JOIN Users u
		  ON u.id = ou.user_id
		  AND u.deleted_dt IS NULL
		INNER JOIN Organization o
		  ON o.id = ou.org_id
		  AND o.deleted_dt IS NULL
		INNER JOIN Meetings m
		  ON ou.org_id = m.org_id
		  AND m.deleted_dt IS NULL
		  AND m.id = ?
		WHERE ou.user_id = ?
		  AND ou.deleted_dt IS NULL`,
		meetingID, userID)
	if err != nil {
		s.logger.Error("GetMeeting Error", err)
		return nil, err
	}
	return &meeting, nil
}

func (s lowerThirdsService) GetMeetings(ctx context.Context) (*[]entities.Meeting, error) {
	userIDVal := ctx.Value(helpers.UserIDKey)
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("userID not found or invalid in context")
	}
	s.logger.Debug("GetMeetings for userID ", userID)

	var meetings []entities.Meeting
	err := s.MySqlDB.Select(
		&meetings,
		`SELECT m.*
		FROM Users u
		INNER JOIN OrgUsers ou
		  ON ou.user_id = u.id
		  AND ou.deleted_dt IS NULL
		INNER JOIN Organization o
		  ON o.id = ou.org_id
		  AND o.deleted_dt IS NULL
		INNER JOIN Meetings m
		  ON ou.org_id = m.org_id
		  AND m.deleted_dt IS NULL
		WHERE u.id = ?
		  AND u.deleted_dt IS NULL`, userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return &meetings, nil
}

func (s lowerThirdsService) UpdateMeeting(ctx context.Context, meetingID uuid.UUID, m *entities.Meeting) error {
	s.logger.Debug("UpdateMeeting")

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.ExecContext(
		ctx,
		`UPDATE Meetings SET
		  id = ?,
		  org_id = ?,
		  conference = ?,
		  meeting = ?,
		  meeting_date = ?,
		  duration = ?
		WHERE id = ?`,
		m.MeetingID,
		m.OrgID,
		m.Conference,
		m.Meeting,
		m.MeetingDate,
		m.Duration,
		meetingID,
	)
	if err != nil {
		s.logger.Error("UpdateMeeting Error", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("UpdateMeeting affected rows: ", affectedRows)
	}
	return nil
}
