package storage

import (
	"lowerthirdsapi/internal/testutil"
	"testing"
)

func TestCreateAndGetUser(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()

	service := New(testutil.TestDB, testutil.TestLogger)

	// Setup test data
	user, _, _ := testutil.CreateTestData(t, service)

	// Test GetUser
	retrievedUser, err := service.GetUser(testutil.TestCtx, user.UserID)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}

	// Verify user data
	if retrievedUser.UserID != user.UserID {
		t.Errorf("Expected UserID %v, got %v", user.UserID, retrievedUser.UserID)
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Expected Email %v, got %v", user.Email, retrievedUser.Email)
	}
	if retrievedUser.FirstName.String != user.FirstName.String {
		t.Errorf("Expected FirstName %v, got %v", user.FirstName.String, retrievedUser.FirstName.String)
	}
	if retrievedUser.LastName.String != user.LastName.String {
		t.Errorf("Expected LastName %v, got %v", user.LastName.String, retrievedUser.LastName.String)
	}
}
