package testutil

import (
	"context"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

// setupTestTables creates all necessary tables for testing
func setupTestTables(db *sqlx.DB) error {
	// Tables are already created in the database, no need to recreate them
	return nil
}

func TestMain(m *testing.M) {
	// Setup test database
	var err error
	// First, connect without specifying a database
	db, err := sqlx.Connect("mysql", "root:@tcp(localhost:3306)/")
	if err != nil {
		panic(err)
	}
	// Create the test database if it doesn't exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS lowerthirds_test")
	if err != nil {
		panic(err)
	}
	db.Close()
	// Now connect to the test database
	TestDB, err = sqlx.Connect("mysql", "root:@tcp(localhost:3306)/lowerthirds_test?parseTime=true")
	if err != nil {
		panic(err)
	}

	// Setup test tables
	if err := setupTestTables(TestDB); err != nil {
		panic(err)
	}

	// Setup test logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	TestLogger = logger.WithField("test", true)

	// Setup test context
	TestCtx = context.Background()
	TestCtx = context.WithValue(TestCtx, helpers.SocialIDKey, "test-social-id")

	// Run tests
	code := m.Run()

	// Cleanup
	TestDB.Close()

	os.Exit(code)
}

func SetupTest(t *testing.T) {
	// Setup test database
	var err error
	// First, connect without specifying a database
	db, err := sqlx.Connect("mysql", "root:@tcp(localhost:3306)/")
	if err != nil {
		t.Fatalf("Failed to connect to MySQL: %v", err)
	}
	// Create the test database if it doesn't exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS lowerthirds_test")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	db.Close()
	// Now connect to the test database
	TestDB, err = sqlx.Connect("mysql", "root:@tcp(localhost:3306)/lowerthirds_test?parseTime=true")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean up any existing test data
	cleanupStmts := []string{
		"DELETE FROM TimerItems WHERE meeting_role = 'Test Role'",
		"DELETE FROM SpeakerItems WHERE meeting_role = 'Test Role'",
		"DELETE FROM MessageItems WHERE meeting_role = 'Test Role'",
		"DELETE FROM LyricsItems WHERE meeting_role = 'Test Role'",
		"DELETE FROM BlankItems WHERE meeting_role = 'Test Role'",
		"DELETE FROM Meetings WHERE meeting = 'Test Meeting'",
		"DELETE FROM OrgUsers WHERE org_id IN (SELECT id FROM Organization WHERE name = 'Test Organization')",
		"DELETE FROM Organization WHERE name = 'Test Organization'",
		"DELETE FROM Users WHERE email = 'test@example.com'",
	}
	for _, stmt := range cleanupStmts {
		_, err = TestDB.Exec(stmt)
		if err != nil {
			t.Fatalf("Failed to clean up test data with statement '%s': %v", stmt, err)
		}
	}

	// Setup test tables
	if err := setupTestTables(TestDB); err != nil {
		t.Fatalf("Failed to setup test tables: %v", err)
	}

	// Setup test logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	TestLogger = logger.WithField("test", true)

	// Setup test context
	TestCtx = context.Background()
	TestCtx = context.WithValue(TestCtx, helpers.SocialIDKey, "test-social-id")
}

func TeardownTest() {
	if TestDB != nil {
		TestDB.Close()
	}
}

// CreateTestData creates test data using the provided service
func CreateTestData(t *testing.T, service interface {
	CreateUser(ctx context.Context, user *entities.User) error
	CreateOrg(ctx context.Context, org *entities.Organization) error
	CreateMeeting(ctx context.Context, meeting *entities.Meeting) error
}) (*entities.User, *entities.Organization, *entities.Meeting) {
	// Create a test user
	userID := uuid.New()
	user := &entities.User{
		UserID:   userID,
		Email:    "test@example.com",
		SocialID: null.StringFrom("test-social-id"),
	}
	err := service.CreateUser(TestCtx, user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create a test organization with a unique name
	orgID := uuid.New()
	org := &entities.Organization{
		OrgID: orgID,
		Name:  "Test Org " + orgID.String(),
	}
	err = service.CreateOrg(TestCtx, org)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}

	// Create a test meeting
	meetingID := uuid.New()
	meeting := &entities.Meeting{
		MeetingID:   meetingID,
		OrgID:       orgID,
		Meeting:     "Test Meeting",
		MeetingDate: time.Now(),
	}
	err = service.CreateMeeting(TestCtx, meeting)
	if err != nil {
		t.Fatalf("Failed to create test meeting: %v", err)
	}

	return user, org, meeting
}
