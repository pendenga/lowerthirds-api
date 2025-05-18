package storage

import (
	"context"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
)

func (s lowerThirdsService) CreateOrg(ctx context.Context, o *entities.Organization) error {
	socialID := ctx.Value(helpers.SocialIDKey).(string)
	s.logger.Debug("CreateOrg for socialID ", socialID)
	user, err := s.GetUserBySocialID(ctx, socialID)
	if err != nil {
		s.logger.Error("User not found by socialID", err)
		return err
	}
	s.logger.Debug("CreateOrg for userID ", user.UserID)

	// TODO: put some user-level security on this query
	_, err = s.MySqlDB.ExecContext(
		ctx,
		`INSERT INTO Organization (id, name) VALUES (?, ?)`,
		o.OrgID,
		o.Name,
	)
	if err != nil {
		s.logger.Error("CreateOrg Error", err)
		return err
	}

	// add permission for this new org
	err = s.CreateOrgUser(ctx, o.OrgID, user.UserID)
	if err != nil {
		s.logger.Error("CreateOrg Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) DeleteOrg(ctx context.Context, orgID uuid.UUID) error {
	s.logger.Debug("DeleteOrg for orgID ", orgID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.ExecContext(ctx, `
		UPDATE Organization 
		SET deleted_dt = CURRENT_TIMESTAMP 
		WHERE id = ?
		  AND deleted_dt IS NULL`,
		orgID,
	)
	if err != nil {
		s.logger.Error("DeleteOrg error ", err)
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("DeleteOrg affected rows: ", affectedRows)

	}
	return nil
}

func (s lowerThirdsService) GetMeetingsByOrg(ctx context.Context, orgID uuid.UUID) (*[]entities.Meeting, error) {
	socialID := ctx.Value(helpers.SocialIDKey).(string)
	s.logger.Debug("GetMeetings for socialID ", socialID)
	user, err := s.GetUserBySocialID(ctx, socialID)
	if err != nil {
		s.logger.Error("User not found by socialID", err)
		return nil, err
	}
	s.logger.Debug("GetMeetings for userID ", user.UserID)

	var meetings []entities.Meeting
	err = s.MySqlDB.Select(
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
		  AND ou.org_id = ?
		  AND ou.deleted_dt IS NULL`, user.UserID, orgID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return &meetings, nil
}

func (s lowerThirdsService) GetOrg(ctx context.Context, orgID uuid.UUID) (*entities.Organization, error) {
	socialID := ctx.Value(helpers.SocialIDKey).(string)
	s.logger.Debug("GetMeeting for socialID ", socialID, " orgID ", orgID)
	user, err := s.GetUserBySocialID(ctx, socialID)
	if err != nil {
		s.logger.Error("User not found by socialID", err)
		return nil, err
	}
	s.logger.Debug("GetMeeting for userID ", user.UserID, " orgID ", orgID)

	var org entities.Organization
	err = s.MySqlDB.Get(
		&org,
		`SELECT o.*
		FROM OrgUsers ou
		INNER JOIN Users u
		  ON u.id = ou.user_id
		  AND u.deleted_dt IS NULL
		INNER JOIN Organization o
		  ON o.id = ou.org_id
		  AND o.deleted_dt IS NULL
		WHERE ou.user_id = ?
		  AND ou.org_id = ?
		  AND ou.deleted_dt IS NULL`,
		user.UserID, orgID)
	if err != nil {
		s.logger.Error("GetOrg Error", err)
		return nil, err
	}
	return &org, nil
}

func (s lowerThirdsService) GetOrgs(ctx context.Context) (*[]entities.Organization, error) {
	socialID := ctx.Value(helpers.SocialIDKey).(string)
	s.logger.Debug("GetMeeting for socialID ", socialID)
	user, err := s.GetUserBySocialID(ctx, socialID)
	if err != nil {
		s.logger.Error("User not found by socialID", err)
		return nil, err
	}
	s.logger.Debug("GetMeeting for userID ", user.UserID)

	var orgs []entities.Organization
	err = s.MySqlDB.Select(
		&orgs,
		`SELECT o.*
		FROM OrgUsers ou
		INNER JOIN Users u
		  ON u.id = ou.user_id
		  AND u.deleted_dt IS NULL
		INNER JOIN Organization o
		  ON o.id = ou.org_id
		  AND o.deleted_dt IS NULL
		WHERE ou.user_id = ?
		  AND ou.deleted_dt IS NULL`,
		user.UserID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return &orgs, nil
}

func (s lowerThirdsService) UpdateOrg(ctx context.Context, orgID uuid.UUID, o *entities.Organization) error {
	s.logger.Debug("UpdateOrg")

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.ExecContext(
		ctx,
		`UPDATE Organization SET id = ?, name = ? WHERE id = ?`,
		o.OrgID,
		o.Name,
		orgID,
	)
	if err != nil {
		s.logger.Error("UpdateOrg Error", err)
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err == nil {
		s.logger.Info("UpdateOrg affected rows: ", affectedRows)
	}
	return nil
}
