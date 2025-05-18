package helpers

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type QueryParams struct {
	Page     int
	PageSize int
	DateFrom *time.Time
	DateTo   *time.Time
	Language string
	UserID   uuid.UUID
	OrgID    uuid.UUID
}

func DefaultQueryParams() QueryParams {
	return QueryParams{
		Page:     0,
		PageSize: 100,
		Language: "eng",
	}
}

func GetQueryParams(ctx context.Context) QueryParams {
	params, ok := ctx.Value(QueryParametersKey).(QueryParams)
	if !ok {
		// Return default values if not set
		return DefaultQueryParams()
	}
	return params
}
