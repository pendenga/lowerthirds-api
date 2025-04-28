package entities

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"time"
)

type Hymn struct {
	HymnID        uuid.UUID   `db:"id" json:"id"`
	Language      string      `db:"language" json:"language"`
	Page          int         `db:"page" json:"page"`
	Name          string      `db:"name" json:"name"`
	TranslationID uuid.UUID   `db:"translation_id" json:"translation_id"`
	Verses        []HymnVerse `json:"verses"`
	DeletedDT     null.Time   `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT    time.Time   `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT     time.Time   `db:"updated_dt" json:"updated_dt"`
}
