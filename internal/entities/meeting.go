package entities

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"time"
)

type Meeting struct {
	MeetingID   uuid.UUID   `db:"id" json:"id"`
	OrgID       uuid.UUID   `db:"org_id" json:"org_id"`
	Conference  null.String `db:"conference" json:"conference"` // nullable STRING
	Meeting     string      `db:"meeting" json:"meeting"`
	MeetingDate time.Time   `db:"meeting_date" json:"date"`
	Duration    null.Int    `db:"duration" json:"duration"` // nullable INT
	AgendaItems []Item      `json:"agenda_items,omitempty"`
	DeletedDT   null.Time   `db:"deleted_dt" json:"deleted_dt"` // nullable DATETIME
	InsertedDT  time.Time   `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT   time.Time   `db:"updated_dt" json:"updated_dt"`
}
