package entities

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"time"
)

type Item interface {
	GetID() uuid.UUID
	GetMeetingID() uuid.UUID
	GetMeetingRole() string
	GetOrder() int
	GetType() string
}

func (b BlankItem) GetID() uuid.UUID   { return b.BlankItemID }
func (m MessageItem) GetID() uuid.UUID { return m.MessageItemID }
func (s SpeakerItem) GetID() uuid.UUID { return s.SpeakerItemID }
func (l LyricsItem) GetID() uuid.UUID  { return l.LyricsItemID }
func (t TimerItem) GetID() uuid.UUID   { return t.TimerItemID }

func (b BlankItem) GetMeetingID() uuid.UUID   { return b.MeetingID }
func (m MessageItem) GetMeetingID() uuid.UUID { return m.MeetingID }
func (s SpeakerItem) GetMeetingID() uuid.UUID { return s.MeetingID }
func (l LyricsItem) GetMeetingID() uuid.UUID  { return l.MeetingID }
func (t TimerItem) GetMeetingID() uuid.UUID   { return t.MeetingID }

func (b BlankItem) GetMeetingRole() string   { return b.MeetingRole }
func (m MessageItem) GetMeetingRole() string { return m.MeetingRole }
func (s SpeakerItem) GetMeetingRole() string { return s.MeetingRole }
func (l LyricsItem) GetMeetingRole() string  { return l.MeetingRole }
func (t TimerItem) GetMeetingRole() string   { return t.MeetingRole }

func (b BlankItem) GetOrder() int   { return b.ItemOrder }
func (m MessageItem) GetOrder() int { return m.ItemOrder }
func (s SpeakerItem) GetOrder() int { return s.ItemOrder }
func (l LyricsItem) GetOrder() int  { return l.ItemOrder }
func (t TimerItem) GetOrder() int   { return t.ItemOrder }

func (b BlankItem) GetType() string   { return b.ItemType }
func (m MessageItem) GetType() string { return m.ItemType }
func (s SpeakerItem) GetType() string { return s.ItemType }
func (l LyricsItem) GetType() string  { return l.ItemType }
func (t TimerItem) GetType() string   { return t.ItemType }

type BlankItem struct {
	BlankItemID uuid.UUID `db:"id" json:"id"`
	MeetingID   uuid.UUID `db:"meeting_id" json:"meeting_id"`
	ItemType    string    `db:"item_type" json:"item_type"`
	ItemOrder   int       `db:"item_order" json:"item_order"`
	MeetingRole string    `db:"meeting_role" json:"meeting_role"`
	DeletedDT   null.Time `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT  time.Time `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT   time.Time `db:"updated_dt" json:"updated_dt"`
}

type LyricsItem struct {
	LyricsItemID    uuid.UUID `db:"id" json:"id"`
	MeetingID       uuid.UUID `db:"meeting_id" json:"meeting_id"`
	ItemType        string    `db:"item_type" json:"item_type"`
	ItemOrder       int       `db:"item_order" json:"item_order"`
	MeetingRole     string    `db:"meeting_role" json:"meeting_role"`
	HymnID          string    `db:"hymn_id" json:"hymn_id"`
	ShowTranslation bool      `db:"show_translation" json:"show_translation"`
	DeletedDT       null.Time `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT      time.Time `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT       time.Time `db:"updated_dt" json:"updated_dt"`
}

type MessageItem struct {
	MessageItemID uuid.UUID   `db:"id" json:"id"`
	MeetingID     uuid.UUID   `db:"meeting_id" json:"meeting_id"`
	ItemType      string      `db:"item_type" json:"item_type"`
	ItemOrder     int         `db:"item_order" json:"item_order"`
	MeetingRole   string      `db:"meeting_role" json:"meeting_role"`
	PrimaryText   string      `db:"primary_text" json:"primary_text"`
	SecondaryText null.String `db:"secondary_text" json:"secondary_text"`
	DeletedDT     null.Time   `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT    time.Time   `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT     time.Time   `db:"updated_dt" json:"updated_dt"`
}

type SpeakerItem struct {
	SpeakerItemID uuid.UUID   `db:"id" json:"id"`
	MeetingID     uuid.UUID   `db:"meeting_id" json:"meeting_id"`
	ItemType      string      `db:"item_type" json:"item_type"`
	ItemOrder     int         `db:"item_order" json:"item_order"`
	MeetingRole   string      `db:"meeting_role" json:"meeting_role"`
	SpeakerName   string      `db:"speaker_name" json:"speaker_name"`
	Title         null.String `db:"title" json:"title"`
	DeletedDT     null.Time   `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT    time.Time   `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT     time.Time   `db:"updated_dt" json:"updated_dt"`
}

type TimerItem struct {
	TimerItemID        uuid.UUID `db:"id" json:"id"`
	MeetingID          uuid.UUID `db:"meeting_id" json:"meeting_id"`
	ItemType           string    `db:"item_type" json:"item_type"`
	ItemOrder          int       `db:"item_order" json:"item_order"`
	MeetingRole        string    `db:"meeting_role" json:"meeting_role"`
	ShowMeetingDetails bool      `db:"show_meeting_details" json:"show_meeting_details"`
	DeletedDT          null.Time `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT         time.Time `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT          time.Time `db:"updated_dt" json:"updated_dt"`
}
