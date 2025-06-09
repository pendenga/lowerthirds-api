package testutil

import (
	"context"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/storage"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v4"
)

var (
	TestDB     *sqlx.DB
	TestLogger *logrus.Entry
	TestCtx    context.Context
)

func SetupTest(t *testing.T) {
	// Setup test database
	var err error
	TestDB, err = sqlx.Connect("mysql", "root:password@tcp(localhost:3306)/lowerthirds_test?parseTime=true")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Setup test logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	TestLogger = logger.WithField("test", true)

	// Setup test context
	TestCtx = context.Background()
	TestCtx = context.WithValue(TestCtx, "socialID", "test-social-id")
}

func TeardownTest() {
	if TestDB != nil {
		TestDB.Close()
	}
}

func SetupTestData(t *testing.T) (*entities.User, *entities.Organization, *entities.Meeting) {
	// Create test user
	userID := uuid.New()
	user := &entities.User{
		UserID:    userID,
		SocialID:  null.StringFrom("test-social-id"),
		Email:     "test@example.com",
		FirstName: null.StringFrom("Test"),
		LastName:  null.StringFrom("User"),
	}

	// Create test organization
	orgID := uuid.New()
	org := &entities.Organization{
		OrgID: orgID,
		Name:  "Test Organization",
	}

	// Create test meeting
	meetingID := uuid.New()
	meeting := &entities.Meeting{
		MeetingID:   meetingID,
		OrgID:       orgID,
		Conference:  null.StringFrom("Test Conference"),
		Meeting:     "Test Meeting",
		MeetingDate: time.Now(),
		Duration:    null.IntFrom(60),
	}

	// Create service instance
	service := storage.New(TestDB, TestLogger)

	// Insert test data
	err := service.CreateUser(TestCtx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	err = service.CreateOrg(TestCtx, org)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}

	err = service.CreateOrgUser(TestCtx, orgID, userID)
	if err != nil {
		t.Fatalf("Failed to create test org user mapping: %v", err)
	}

	err = service.CreateMeeting(TestCtx, meeting)
	if err != nil {
		t.Fatalf("Failed to create test meeting: %v", err)
	}

	return user, org, meeting
}
