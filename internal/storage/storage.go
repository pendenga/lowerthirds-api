package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"lowerthirdsapi/internal/entities"
)

type LowerThirdsService interface {
	// Meetings
	CreateMeeting(ctx context.Context, m *entities.Meeting) error
	DeleteMeeting(ctx context.Context, meetingID uuid.UUID) error
	GetMeeting(ctx context.Context, meetingID uuid.UUID) (*entities.Meeting, error)
	GetMeetings(ctx context.Context) (*[]entities.Meeting, error)
	GetMeetingsByOrg(ctx context.Context, orgID uuid.UUID) (*[]entities.Meeting, error)
	GetMeetingsByUser(ctx context.Context, userID uuid.UUID) (*[]entities.Meeting, error)
	UpdateMeeting(ctx context.Context, meetingID uuid.UUID, m *entities.Meeting) error

	// Orgs
	CreateOrg(ctx context.Context, m *entities.Organization) error
	DeleteOrg(ctx context.Context, orgID uuid.UUID) error
	GetOrg(ctx context.Context, orgID uuid.UUID) (*entities.Organization, error)
	GetOrgs(ctx context.Context) (*[]entities.Organization, error)
	UpdateOrg(ctx context.Context, orgID uuid.UUID, m *entities.Organization) error

	// OrgUser
	CreateOrgUser(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error
	DeleteOrgUser(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error
	GetOrgsByUser(ctx context.Context, userID uuid.UUID) (*[]entities.Organization, error)
	SetOrgsByUser(ctx context.Context, userID uuid.UUID, orgIDs []uuid.UUID) error
	GetUsersByOrg(ctx context.Context, orgID uuid.UUID) (*[]entities.User, error)

	// Slides
	CreateSlide(ctx context.Context, slide entities.Slide) error
	DeleteSlide(ctx context.Context, slideID uuid.UUID) error
	GetSlide(ctx context.Context, slideID uuid.UUID) (entities.Slide, error)
	GetSlides(ctx context.Context) (*[]entities.Slide, error)
	GetSlidesByMeeting(ctx context.Context, meetingID uuid.UUID) (*[]entities.Slide, error)
	UpdateSlide(ctx context.Context, slideID uuid.UUID, slide entities.Slide) error

	// Users
	CreateUser(ctx context.Context, m *entities.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	GetUser(ctx context.Context, userID uuid.UUID) (*entities.User, error)
	GetUsers(ctx context.Context) (*[]entities.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, m *entities.User) error
}

type lowerThirdsService struct {
	MySqlDB *sqlx.DB
	logger  *logrus.Entry
}

func New(db *sqlx.DB, l *logrus.Entry) LowerThirdsService {
	return &lowerThirdsService{
		MySqlDB: db,
		logger:  l,
	}
}
