package storage

import (
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/storage/testutil"
	"testing"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

func TestCreateAndGetMeeting(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()

	service := New(testutil.TestDB, testutil.TestLogger)

	// Create test meeting
	meetingID := uuid.New()
	orgID := uuid.New()
	meeting := &entities.Meeting{
		MeetingID:   meetingID,
		OrgID:       orgID,
		Conference:  null.StringFrom("Test Conference 2"),
		Meeting:     "Test Meeting 2",
		MeetingDate: time.Now(),
		Duration:    null.IntFrom(90),
	}

	// Test CreateMeeting
	err := service.CreateMeeting(testutil.TestCtx, meeting)
	if err != nil {
		t.Fatalf("CreateMeeting failed: %v", err)
	}

	// Test GetMeeting
	retrievedMeeting, err := service.GetMeeting(testutil.TestCtx, meetingID)
	if err != nil {
		t.Fatalf("GetMeeting failed: %v", err)
	}

	// Verify meeting data
	if retrievedMeeting.MeetingID != meeting.MeetingID {
		t.Errorf("Expected MeetingID %v, got %v", meeting.MeetingID, retrievedMeeting.MeetingID)
	}
	if retrievedMeeting.OrgID != meeting.OrgID {
		t.Errorf("Expected OrgID %v, got %v", meeting.OrgID, retrievedMeeting.OrgID)
	}
	if retrievedMeeting.Conference.String != meeting.Conference.String {
		t.Errorf("Expected Conference %v, got %v", meeting.Conference.String, retrievedMeeting.Conference.String)
	}
	if retrievedMeeting.Meeting != meeting.Meeting {
		t.Errorf("Expected Meeting %v, got %v", meeting.Meeting, retrievedMeeting.Meeting)
	}
	if retrievedMeeting.Duration.Int64 != meeting.Duration.Int64 {
		t.Errorf("Expected Duration %v, got %v", meeting.Duration.Int64, retrievedMeeting.Duration.Int64)
	}
}
