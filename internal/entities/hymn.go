package entities

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"time"
)

type Hymn struct {
	HymnID     uuid.UUID   `db:"id" json:"id"`
	Language   string      `db:"language" json:"language"`
	Page       null.String `db:"page" json:"page"`
	DeletedDT  null.Time   `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT time.Time   `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT  time.Time   `db:"updated_dt" json:"updated_dt"`
}
