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
    GetMeetingsByUser(ctx context.Context, userID string) (*[]entities.Meeting, error)

    // Orgs
    CreateOrg(ctx context.Context, o *entities.Organization) error
    DeleteOrg(ctx context.Context, orgID uuid.UUID) error
    GetOrg(ctx context.Context, orgID uuid.UUID) (*entities.Organization, error)
    GetOrgs(ctx context.Context) (*[]entities.Organization, error)
    UpdateOrg(ctx context.Context, orgID uuid.UUID, o *entities.Organization) error

    // OrgUser
    CreateOrgUser(ctx context.Context, orgID uuid.UUID, userID string) error
    DeleteOrgUser(ctx context.Context, orgID uuid.UUID, userID string) error
    GetOrgsByUser(ctx context.Context, userID string) (*[]entities.Organization, error)
    SetOrgsByUser(ctx context.Context, userID string, orgIDs []uuid.UUID) error
    GetUsersByOrg(ctx context.Context, orgID uuid.UUID) (*[]entities.User, error)

    // Slides
    CreateSlide(ctx context.Context, slide entities.Slide) error
    DeleteSlide(ctx context.Context, slideID uuid.UUID) error
    GetSlide(ctx context.Context, slideID uuid.UUID) (entities.Slide, error)
    GetSlides(ctx context.Context) (*[]entities.Slide, error)
    GetSlidesByMeeting(ctx context.Context, meetingID uuid.UUID) (*[]entities.Slide, error)
    UpdateSlide(ctx context.Context, slideID uuid.UUID, slide entities.Slide) error

    // Users
    CreateUser(ctx context.Context, u *entities.User) error
    DeleteUser(ctx context.Context, userID string) error
    GetUser(ctx context.Context, userID string) (*entities.User, error)
    GetUsers(ctx context.Context) (*[]entities.User, error)
    UpdateUser(ctx context.Context, userID string, u *entities.User) error
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
