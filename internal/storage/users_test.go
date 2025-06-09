package storage

import (
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/storage/testutil"
	"testing"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

func TestCreateAndGetUser(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()

	service := New(testutil.TestDB, testutil.TestLogger)

	// Create test user
	userID := uuid.New()
	user := &entities.User{
		UserID:    userID,
		SocialID:  null.StringFrom("test-social-id-2"),
		Email:     "test2@example.com",
		FirstName: null.StringFrom("Test2"),
		LastName:  null.StringFrom("User2"),
	}

	// Test CreateUser
	err := service.CreateUser(testutil.TestCtx, user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Test GetUser
	retrievedUser, err := service.GetUser(testutil.TestCtx, userID)
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
