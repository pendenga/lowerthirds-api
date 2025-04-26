package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
)

func (s lowerThirdsService) CreateOrg(ctx context.Context, o *entities.Organization) error {
	userIDVal := ctx.Value(helpers.UserIDKey)
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		return fmt.Errorf("userID not found or invalid in context")
	}
	s.logger.Debug("CreateOrg for userID ", userID)

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.ExecContext(
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
	err = s.CreateOrgUser(ctx, o.OrgID, userID)
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
		  AND ou.deleted_dt IS NULL`, userID, orgID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return &meetings, nil
}

func (s lowerThirdsService) GetOrg(ctx context.Context, orgID uuid.UUID) (*entities.Organization, error) {
	userIDVal := ctx.Value(helpers.UserIDKey)
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("userID not found or invalid in context")
	}
	s.logger.Debug("GetOrg for userID ", userID, " orgID ", orgID)
	var org entities.Organization
	err := s.MySqlDB.Get(
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
		userID, orgID)
	if err != nil {
		s.logger.Error("GetOrg Error", err)
		return nil, err
	}
	return &org, nil
}

func (s lowerThirdsService) GetOrgs(ctx context.Context) (*[]entities.Organization, error) {
	userIDVal := ctx.Value(helpers.UserIDKey)
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("userID not found or invalid in context")
	}
	s.logger.Debug("GetOrgs for userID ", userID)

	var orgs []entities.Organization
	err := s.MySqlDB.Select(
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
		userID)
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
