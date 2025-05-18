package storage

import (
	"context"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
)

func (s lowerThirdsService) CreateUser(ctx context.Context, u *entities.User) error {
	s.logger.Debug("CreateUser")

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.ExecContext(
		ctx,
		`INSERT INTO Users (id, email, first_name, full_name, last_name, social_id, photo_url) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		u.UserID,
		u.Email,
		u.FirstName,
		u.FullName,
		u.LastName,
		u.SocialID,
		u.PhotoURL,
	)
	if err != nil {
		s.logger.Error("CreateUser Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	s.logger.Debug("DeleteUser for userID ", userID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.ExecContext(ctx, `
		UPDATE Users 
		SET deleted_dt = CURRENT_TIMESTAMP 
		WHERE id = ?
		  AND deleted_dt IS NULL`,
		userID,
	)
	if err != nil {
		s.logger.Error("DeleteUser error ", err)
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("DeleteUser affected rows: ", affectedRows)

	}
	return nil
}

func (s lowerThirdsService) GetMeetingsByUser(ctx context.Context, userID uuid.UUID) (*[]entities.Meeting, error) {
	s.logger.Debug("GetMeetings for userID ", userID)

	var meetings []entities.Meeting
	err := s.MySqlDB.Select(
		&meetings,
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
		WHERE ou.user_id = ?
		  AND ou.deleted_dt IS NULL`, userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return &meetings, nil
}

func (s lowerThirdsService) GetUser(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	s.logger.Debug("GetUser for userID ", userID)
	var user entities.User
	err := s.MySqlDB.Get(&user, `SELECT * FROM Users WHERE id = ?`, userID)
	if err != nil {
		s.logger.Error("GetUser Error", err)
		return nil, err
	}
	return &user, nil
}

func (s lowerThirdsService) GetUsers(ctx context.Context) (*[]entities.User, error) {
	s.logger.Debug("GetUsers")

	var users []entities.User
	err := s.MySqlDB.Select(&users, `SELECT * FROM Users`)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return &users, nil
}

func (s lowerThirdsService) UpdateUser(ctx context.Context, userID uuid.UUID, u *entities.User) error {
	s.logger.Debug("UpdateUser")

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.ExecContext(
		ctx,
		`UPDATE Users 
        SET id = ?, 
          email = ?, 
          first_name = ?, 
          full_name = ?, 
          last_name = ?
          social_id = ?
          photo_url = ?
        WHERE id = ?`,
		u.UserID,
		u.Email,
		u.FirstName,
		u.FullName,
		u.LastName,
		u.SocialID,
		u.PhotoURL,
		userID,
	)
	if err != nil {
		s.logger.Error("UpdateUser Error", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("UpdateUser affected rows: ", affectedRows)
	}
	return nil
}
