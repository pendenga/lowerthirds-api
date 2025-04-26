package entities

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"time"
)

type OrgUser struct {
	OrgID      uuid.UUID `db:"org_id" json:"org_id"`
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	DeletedDT  null.Time `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT time.Time `db:"inserted_dt" json:"inserted_dt"`
}
