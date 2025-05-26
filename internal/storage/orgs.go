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

	// check if user.UserID is already in o.UserIDs
	var nu []uuid.UUID
	found := false
	for _, userID := range o.UserIDs {
		nu = append(nu, userID)
		if userID == user.UserID {
			found = true
		}
	}
	if !found {
		nu = append(nu, user.UserID)
	}

	// After updating the org, we need to validate the user list
	var ex []uuid.UUID
	affectedRows, err := s.reconcileUsers(ctx, o.OrgID, ex, nu)
	if err == nil {
		s.logger.Info("CreateOrg affected users: ", affectedRows)
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

	// delete all users
	ex, err := s.GetUserIDsByOrg(ctx, orgID)
	if err != nil {
		s.logger.Error("DeleteOrg Users Error", err)
		return err
	}

	// After updating the org, we need to validate the user list
	var nu []uuid.UUID
	affectedRows, err = s.reconcileUsers(ctx, orgID, *ex, nu)
	if err == nil {
		s.logger.Info("DeleteOrg affected users: ", affectedRows)
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

	// assign all users to all orgs
	userIDs, err := s.GetUserIDsByOrg(ctx, orgID)
	if err != nil {
		s.logger.Error("GetOrg Users Error", err)
		return nil, err
	}
	org.UserIDs = *userIDs

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

	// assign all users to all orgs
	OrgUsersMap, err := s.GetOrgUsersMap(ctx)
	if err != nil {
		s.logger.Error("GetOrgs Users Error", err)
		return nil, err
	}

	var returnOrgs []entities.Organization
	for _, org := range orgs {
		org.UserIDs = OrgUsersMap[org.OrgID]
		returnOrgs = append(returnOrgs, org)
	}

	return &returnOrgs, nil
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

	// get existing users for the org
	existingUserIDs, err := s.GetUserIDsByOrg(ctx, orgID)
	if err != nil {
		s.logger.Error("GetUsersByOrg Error", err)
		return err
	}

	// After updating the org, we need to validate the user list
	affectedRows, err = s.reconcileUsers(ctx, orgID, *existingUserIDs, o.UserIDs)
	if err == nil {
		s.logger.Info("UpdateOrg affected rows: ", affectedRows)
	}

	return nil
}

func (s lowerThirdsService) reconcileUsers(ctx context.Context, orgID uuid.UUID, ex []uuid.UUID, nu []uuid.UUID) (int64, error) {
	// make a map of existing users
	existingUserMap := make(map[uuid.UUID]bool)
	for _, userID := range ex {
		existingUserMap[userID] = true
	}

	var affected int64 = 0
	newUserMap := make(map[uuid.UUID]bool)
	for _, userID := range nu {
		// Add to map of new users
		newUserMap[userID] = true

		// If the user is in the new map but not in the existing map, add them
		if _, exists := existingUserMap[userID]; !exists {
			err := s.CreateOrgUser(ctx, orgID, userID)
			if err != nil {
				s.logger.Error("CreateOrgUser Error", err)
				return 0, err
			}
			affected++
		}
	}

	// If the user is in the existing map but not in the new map, remove them
	for _, userID := range ex {
		if _, exists := newUserMap[userID]; !exists {
			err := s.DeleteOrgUser(ctx, orgID, userID)
			if err != nil {
				s.logger.Error("DeleteOrgUser Error", err)
				return 0, err
			}
			affected++
		}
	}

	return affected, nil
}
