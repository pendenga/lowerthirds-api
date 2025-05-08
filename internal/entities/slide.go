package entities

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"time"
)

type Slide interface {
	GetID() uuid.UUID
	GetMeetingID() uuid.UUID
	GetMeetingRole() string
	GetOrder() int
	GetType() string
}

func (b BlankSlide) GetID() uuid.UUID   { return b.BlankSlideID }
func (m MessageSlide) GetID() uuid.UUID { return m.MessageSlideID }
func (s SpeakerSlide) GetID() uuid.UUID { return s.SpeakerSlideID }
func (l LyricsSlide) GetID() uuid.UUID  { return l.LyricsSlideID }
func (t TimerSlide) GetID() uuid.UUID   { return t.TimerSlideID }

func (b BlankSlide) GetMeetingID() uuid.UUID   { return b.MeetingID }
func (m MessageSlide) GetMeetingID() uuid.UUID { return m.MeetingID }
func (s SpeakerSlide) GetMeetingID() uuid.UUID { return s.MeetingID }
func (l LyricsSlide) GetMeetingID() uuid.UUID  { return l.MeetingID }
func (t TimerSlide) GetMeetingID() uuid.UUID   { return t.MeetingID }

func (b BlankSlide) GetMeetingRole() string   { return b.MeetingRole }
func (m MessageSlide) GetMeetingRole() string { return m.MeetingRole }
func (s SpeakerSlide) GetMeetingRole() string { return s.MeetingRole }
func (l LyricsSlide) GetMeetingRole() string  { return l.MeetingRole }
func (t TimerSlide) GetMeetingRole() string   { return t.MeetingRole }

func (b BlankSlide) GetOrder() int   { return b.SlideOrder }
func (m MessageSlide) GetOrder() int { return m.SlideOrder }
func (s SpeakerSlide) GetOrder() int { return s.SlideOrder }
func (l LyricsSlide) GetOrder() int  { return l.SlideOrder }
func (t TimerSlide) GetOrder() int   { return t.SlideOrder }

func (b BlankSlide) GetType() string   { return b.SlideType }
func (m MessageSlide) GetType() string { return m.SlideType }
func (s SpeakerSlide) GetType() string { return s.SlideType }
func (l LyricsSlide) GetType() string  { return l.SlideType }
func (t TimerSlide) GetType() string   { return t.SlideType }

type BlankSlide struct {
	BlankSlideID uuid.UUID `db:"id" json:"id"`
	MeetingID    uuid.UUID `db:"meeting_id" json:"meeting_id"`
	SlideType    string    `db:"slide_type" json:"slide_type"`
	SlideOrder   int       `db:"slide_order" json:"slide_order"`
	MeetingRole  string    `db:"meeting_role" json:"meeting_role"`
	DeletedDT    null.Time `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT   time.Time `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT    time.Time `db:"updated_dt" json:"updated_dt"`
}

type LyricsSlide struct {
	LyricsSlideID   uuid.UUID `db:"id" json:"id"`
	MeetingID       uuid.UUID `db:"meeting_id" json:"meeting_id"`
	SlideType       string    `db:"slide_type" json:"slide_type"`
	SlideOrder      int       `db:"slide_order" json:"slide_order"`
	MeetingRole     string    `db:"meeting_role" json:"meeting_role"`
	HymnID          string    `db:"hymn_id" json:"hymn_id"`
	ShowTranslation bool      `db:"show_translation" json:"show_translation"`
	DeletedDT       null.Time `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT      time.Time `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT       time.Time `db:"updated_dt" json:"updated_dt"`
}

type MessageSlide struct {
	MessageSlideID uuid.UUID   `db:"id" json:"id"`
	MeetingID      uuid.UUID   `db:"meeting_id" json:"meeting_id"`
	SlideType      string      `db:"slide_type" json:"slide_type"`
	SlideOrder     int         `db:"slide_order" json:"slide_order"`
	MeetingRole    string      `db:"meeting_role" json:"meeting_role"`
	PrimaryText    string      `db:"primary_text" json:"primary_text"`
	SecondaryText  null.String `db:"secondary_text" json:"secondary_text"`
	DeletedDT      null.Time   `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT     time.Time   `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT      time.Time   `db:"updated_dt" json:"updated_dt"`
}

type SpeakerSlide struct {
	SpeakerSlideID uuid.UUID   `db:"id" json:"id"`
	MeetingID      uuid.UUID   `db:"meeting_id" json:"meeting_id"`
	SlideType      string      `db:"slide_type" json:"slide_type"`
	SlideOrder     int         `db:"slide_order" json:"slide_order"`
	MeetingRole    string      `db:"meeting_role" json:"meeting_role"`
	SpeakerName    string      `db:"speaker_name" json:"speaker_name"`
	Title          null.String `db:"title" json:"title"`
	DeletedDT      null.Time   `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT     time.Time   `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT      time.Time   `db:"updated_dt" json:"updated_dt"`
}

type TimerSlide struct {
	TimerSlideID       uuid.UUID `db:"id" json:"id"`
	MeetingID          uuid.UUID `db:"meeting_id" json:"meeting_id"`
	SlideType          string    `db:"slide_type" json:"slide_type"`
	SlideOrder         int       `db:"slide_order" json:"slide_order"`
	MeetingRole        string    `db:"meeting_role" json:"meeting_role"`
	ShowMeetingDetails bool      `db:"show_meeting_details" json:"show_meeting_details"`
	DeletedDT          null.Time `db:"deleted_dt" json:"deleted_dt"`
	InsertedDT         time.Time `db:"inserted_dt" json:"inserted_dt"`
	UpdatedDT          time.Time `db:"updated_dt" json:"updated_dt"`
}
