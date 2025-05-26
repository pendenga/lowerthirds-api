package entities

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"time"
)

type Organization struct {
	OrgID      uuid.UUID   `db:"id" json:"id"`
	Name       string      `db:"name" json:"name"`
	UserIDs    []uuid.UUID `json:"user_ids"`
	DeletedDT  null.Time   `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT time.Time   `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT  time.Time   `db:"updated_dt" json:"updated_dt"`
}
