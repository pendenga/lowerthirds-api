package storage

import (
	"context"
	"lowerthirdsapi/internal/entities"
)

func (s lowerThirdsService) GetUserBySocialID(ctx context.Context, socialID string) (*entities.User, error) {
	// TODO: setting this in ctx...
	// userID := ctx.Value(helpers.UserIDKey).(string)
	s.logger.Debug("GetUserBySocialID for socialID ", socialID)
	var user entities.User
	err := s.MySqlDB.Get(
		&user,
		`SELECT *
        FROM Users
        WHERE social_id = ?
          AND deleted_dt IS NULL`,
		socialID,
	)
	if err != nil {
		s.logger.Error("GetUserBySocialID Error", err)
		return nil, err
	}
	return &user, nil
}
