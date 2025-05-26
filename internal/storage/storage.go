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
	UpdateMeeting(ctx context.Context, meetingID uuid.UUID, m *entities.Meeting) error
	GetMeetingsByOrg(ctx context.Context, orgID uuid.UUID) (*[]entities.Meeting, error)
	GetMeetingsByUser(ctx context.Context, userID uuid.UUID) (*[]entities.Meeting, error)

	// Orgs
	CreateOrg(ctx context.Context, o *entities.Organization) error
	DeleteOrg(ctx context.Context, orgID uuid.UUID) error
	GetOrg(ctx context.Context, orgID uuid.UUID) (*entities.Organization, error)
	GetOrgs(ctx context.Context) (*[]entities.Organization, error)
	UpdateOrg(ctx context.Context, orgID uuid.UUID, o *entities.Organization) error

	// OrgUser
	CreateOrgUser(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error
	DeleteOrgUser(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error
	GetOrgsByUser(ctx context.Context, userID uuid.UUID) (*[]entities.Organization, error)
	SetOrgsByUser(ctx context.Context, userID uuid.UUID, orgIDs []uuid.UUID) error
	GetOrgUsersMap(ctx context.Context) (map[uuid.UUID][]uuid.UUID, error)
	GetUsersByOrg(ctx context.Context, orgID uuid.UUID) (*[]entities.User, error)

	// Items
	CreateItem(ctx context.Context, item entities.Item) error
	DeleteItem(ctx context.Context, itemID uuid.UUID) error
	GetItem(ctx context.Context, itemID uuid.UUID) (entities.Item, error)
	GetItems(ctx context.Context) (*[]entities.Item, error)
	GetItemsByMeeting(ctx context.Context, meetingID uuid.UUID) (*[]entities.Item, error)
	UpdateItem(ctx context.Context, itemID uuid.UUID, item entities.Item) error

	// Users
	CreateUser(ctx context.Context, u *entities.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	GetUser(ctx context.Context, userID uuid.UUID) (*entities.User, error)
	GetUsers(ctx context.Context) (*[]entities.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, u *entities.User) error
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
