package storage

import (
	"context"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
)

func (s lowerThirdsService) CreateOrgUser(ctx context.Context, orgID uuid.UUID, userID string) error {
	s.logger.Debug("CreateOrgUser")

	// TODO: put some user-level security on this query
	_, err := s.MySqlDB.ExecContext(
		ctx,
		`INSERT INTO OrgUsers (org_id, user_id) VALUES (?, ?)`,
		orgID,
		userID,
	)
	if err != nil {
		s.logger.Error("CreateOrgUser Error", err)
		return err
	}
	return nil
}

func (s lowerThirdsService) DeleteOrgUser(ctx context.Context, orgID uuid.UUID, userID string) error {
	s.logger.Debug("DeleteOrg for orgID ", orgID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.ExecContext(ctx, `
		UPDATE OrgUsers 
		SET deleted_dt = CURRENT_TIMESTAMP 
		WHERE org_id = ? 
		  AND user_id = ?
		  AND deleted_dt IS NULL`,
		orgID,
		userID,
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

func (s lowerThirdsService) DeleteOrgsByUser(ctx context.Context, userID string) error {
	s.logger.Debug("DeleteOrgs for userID ", userID)

	// TODO: put some user level security on this query
	result, err := s.MySqlDB.ExecContext(ctx, `
		UPDATE OrgUsers 
		SET deleted_dt = CURRENT_TIMESTAMP 
		WHERE user_id = ?
		  AND deleted_dt IS NULL`,
		userID,
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

func (s lowerThirdsService) GetUsersByOrg(ctx context.Context, orgID uuid.UUID) (*[]entities.User, error) {
	s.logger.Debug("GetUsers for orgID ", orgID)

	var users []entities.User
	err := s.MySqlDB.Select(
		&users,
		`SELECT u.*
		FROM OrgUsers ou
		INNER JOIN Users u
		  ON u.id = ou.user_id
		  AND u.deleted_dt IS NULL
		INNER JOIN Organization o
		  ON o.id = ou.org_id
		  AND o.deleted_dt IS NULL
		WHERE ou.org_id = ?
		  AND ou.deleted_dt IS NULL`,
		orgID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return &users, nil
}

func (s lowerThirdsService) GetOrgsByUser(ctx context.Context, userID string) (*[]entities.Organization, error) {
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

func (s lowerThirdsService) SetOrgsByUser(ctx context.Context, userID string, orgIDs []uuid.UUID) error {
	s.logger.Debug("SetOrgsByUser for userID ", userID, ", orgIDs ", orgIDs)

	err := s.DeleteOrgsByUser(ctx, userID)
	if err != nil {
		s.logger.Error("SetOrgsByUser delete error", err)
		return err
	}

	for _, orgID := range orgIDs {
		err := s.CreateOrgUser(ctx, orgID, userID)
		if err != nil {
			s.logger.Error("SetOrgsByUser error ", err)
			return err
		}
	}
	return nil
}
