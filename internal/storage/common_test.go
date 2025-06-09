package storage

import (
	"context"
	"lowerthirdsapi/internal/entities"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v4"
)

var (
	testDB     *sqlx.DB
	testLogger *logrus.Entry
	testCtx    context.Context
)

func TestMain(m *testing.M) {
	// Setup test database
	var err error
	testDB, err = sqlx.Connect("mysql", "root:password@tcp(localhost:3306)/lowerthirds_test?parseTime=true")
	if err != nil {
		panic(err)
	}

	// Setup test logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	testLogger = logger.WithField("test", true)

	// Setup test context
	testCtx = context.Background()
	testCtx = context.WithValue(testCtx, "socialID", "test-social-id")

	// Run tests
	code := m.Run()

	// Cleanup
	testDB.Close()

	os.Exit(code)
}

func setupTestData(t *testing.T) (*entities.User, *entities.Organization, *entities.Meeting) {
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
	service := New(testDB, testLogger)

	// Insert test data
	err := service.CreateUser(testCtx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	err = service.CreateOrg(testCtx, org)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}

	err = service.CreateOrgUser(testCtx, orgID, userID)
	if err != nil {
		t.Fatalf("Failed to create test org user mapping: %v", err)
	}

	err = service.CreateMeeting(testCtx, meeting)
	if err != nil {
		t.Fatalf("Failed to create test meeting: %v", err)
	}

	return user, org, meeting
}
