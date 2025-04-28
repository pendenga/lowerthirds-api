package entities

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"time"
)

type HymnVerse struct {
	HymnID      uuid.UUID   `db:"hymn_id" json:"hymn_id"`
	VerseNumber int         `db:"verse_number" json:"number"`
	VerseLines  null.String `db:"verse_lines" json:"lines"`
	Optional    bool        `db:"optional" json:"optional"`
	DeletedDT   null.Time   `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT  time.Time   `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT   time.Time   `db:"updated_dt" json:"updated_dt"`
}
