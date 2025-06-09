package storage

import (
	"context"
	"fmt"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/testutil"
	"strings"
	"testing"
	"time"

	"lowerthirdsapi/internal/helpers"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

func TestCreateAndGetItem(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()

	service := New(testutil.TestDB, testutil.TestLogger)

	// Setup test data
	_, _, meeting := testutil.CreateTestData(t, service)

	// Create test message item
	messageItemID := uuid.New()
	messageItem := &entities.MessageItem{
		MessageItemID: messageItemID,
		MeetingID:     meeting.MeetingID,
		ItemType:      "message",
		ItemOrder:     1,
		MeetingRole:   "Test Role",
		PrimaryText:   "Test Message",
	}

	// Create the item
	err := service.CreateItem(testutil.TestCtx, messageItem)
	if err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}

	// Get the item
	retrievedItem, err := service.GetItem(testutil.TestCtx, messageItemID)
	if err != nil {
		t.Fatalf("GetItem failed: %v", err)
	}

	// Verify the item
	retrievedMessageItem, ok := retrievedItem.(*entities.MessageItem)
	if !ok {
		t.Fatalf("Expected MessageItem, got %T", retrievedItem)
	}

	if retrievedMessageItem.MessageItemID != messageItemID {
		t.Errorf("Expected MessageItemID %v, got %v", messageItemID, retrievedMessageItem.MessageItemID)
	}
	if retrievedMessageItem.MeetingID != meeting.MeetingID {
		t.Errorf("Expected MeetingID %v, got %v", meeting.MeetingID, retrievedMessageItem.MeetingID)
	}
	if retrievedMessageItem.ItemType != "message" {
		t.Errorf("Expected ItemType 'message', got %v", retrievedMessageItem.ItemType)
	}
	if retrievedMessageItem.ItemOrder != 1 {
		t.Errorf("Expected ItemOrder 1, got %v", retrievedMessageItem.ItemOrder)
	}
	if retrievedMessageItem.MeetingRole != "Test Role" {
		t.Errorf("Expected MeetingRole 'Test Role', got %v", retrievedMessageItem.MeetingRole)
	}
	if retrievedMessageItem.PrimaryText != "Test Message" {
		t.Errorf("Expected PrimaryText 'Test Message', got %v", retrievedMessageItem.PrimaryText)
	}
}

func TestDeleteItem(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()

	service := New(testutil.TestDB, testutil.TestLogger)

	// Setup test data
	_, _, meeting := testutil.CreateTestData(t, service)

	// Create test message item
	messageItemID := uuid.New()
	messageItem := &entities.MessageItem{
		MessageItemID: messageItemID,
		MeetingID:     meeting.MeetingID,
		ItemType:      "message",
		ItemOrder:     1,
		MeetingRole:   "Test Role",
		PrimaryText:   "Test Message",
	}

	// Create the item
	err := service.CreateItem(testutil.TestCtx, messageItem)
	if err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}

	// Delete the item
	err = service.DeleteItem(testutil.TestCtx, messageItemID)
	if err != nil {
		t.Fatalf("DeleteItem failed: %v", err)
	}

	// Verify the item is deleted
	_, err = service.GetItem(testutil.TestCtx, messageItemID)
	if err == nil {
		t.Fatalf("Expected error when getting deleted item, got nil")
	}
	if err.Error() != "item not found" {
		t.Fatalf("Expected 'item not found' error, got: %v", err)
	}
}

func TestUpdateItem(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()

	service := New(testutil.TestDB, testutil.TestLogger)

	// Setup test data
	_, _, meeting := testutil.CreateTestData(t, service)
	itemID := uuid.New()
	blankItem := &entities.BlankItem{
		BlankItemID: itemID,
		MeetingID:   meeting.MeetingID,
		ItemType:    "blank",
		ItemOrder:   1,
		MeetingRole: "Test Role",
	}
	err := service.CreateItem(testutil.TestCtx, blankItem)
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	// Update the item
	blankItem.ItemOrder = 2
	err = service.UpdateItem(testutil.TestCtx, itemID, blankItem)
	if err != nil {
		t.Fatalf("Failed to update item: %v", err)
	}

	// Verify the update
	updatedItem, err := service.GetItem(testutil.TestCtx, itemID)
	if err != nil {
		t.Fatalf("Failed to get updated item: %v", err)
	}
	updatedBlank, ok := updatedItem.(*entities.BlankItem)
	if !ok {
		t.Fatalf("Expected *entities.BlankItem, got %T", updatedItem)
	}
	if updatedBlank.ItemOrder != 2 {
		t.Errorf("Expected ItemOrder to be 2, got %d", updatedBlank.ItemOrder)
	}
}

func TestGetItemsByMeeting(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)

	// Setup test data
	_, _, meeting := testutil.CreateTestData(t, service)
	itemID := uuid.New()
	blankItem := &entities.BlankItem{
		BlankItemID: itemID,
		MeetingID:   meeting.MeetingID,
		ItemType:    "blank",
		ItemOrder:   1,
		MeetingRole: "Test Role",
	}
	err := service.CreateItem(testutil.TestCtx, blankItem)
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	// Get items by meeting
	items, err := service.GetItemsByMeeting(testutil.TestCtx, meeting.MeetingID)
	if err != nil {
		t.Fatalf("Failed to get items by meeting: %v", err)
	}
	if items == nil {
		t.Fatalf("Expected non-nil items slice")
	}
	if len(*items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(*items))
	}

	// Debug: Print all items
	for i, item := range *items {
		t.Logf("Item %d: Type=%T, ID=%v", i, item, item.GetID())
	}

	// Find our test item in the list
	var foundBlankItem bool
	for _, item := range *items {
		blank, ok := item.(*entities.BlankItem)
		if ok && blank.BlankItemID == itemID {
			foundBlankItem = true
			if blank.MeetingID != meeting.MeetingID {
				t.Errorf("Expected MeetingID %v, got %v", meeting.MeetingID, blank.MeetingID)
			}
			if blank.ItemType != "blank" {
				t.Errorf("Expected ItemType 'blank', got %v", blank.ItemType)
			}
			if blank.ItemOrder != 1 {
				t.Errorf("Expected ItemOrder 1, got %v", blank.ItemOrder)
			}
			if blank.MeetingRole != "Test Role" {
				t.Errorf("Expected MeetingRole 'Test Role', got %v", blank.MeetingRole)
			}
			break
		}
	}
	if !foundBlankItem {
		t.Fatalf("Expected to find BlankItem with ID %v", itemID)
	}
}

func TestGetItems(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)

	// Clean up all relevant tables
	cleanupStmts := []string{
		"DELETE FROM TimerItems",
		"DELETE FROM SpeakerItems",
		"DELETE FROM MessageItems",
		"DELETE FROM LyricsItems",
		"DELETE FROM BlankItems",
		"DELETE FROM Meetings",
		"DELETE FROM OrgUsers",
		"DELETE FROM Organization",
		"DELETE FROM Users",
	}
	for _, stmt := range cleanupStmts {
		_, err := testutil.TestDB.Exec(stmt)
		if err != nil {
			t.Fatalf("Failed to clean up test data: %v", err)
		}
	}

	// Setup test data
	_, _, meeting := testutil.CreateTestData(t, service)

	// Create test items
	blankItemID := uuid.New()
	blankItem := &entities.BlankItem{
		BlankItemID: blankItemID,
		MeetingID:   meeting.MeetingID,
		ItemType:    "blank",
		ItemOrder:   1,
		MeetingRole: "Test Role",
	}
	err := service.CreateItem(testutil.TestCtx, blankItem)
	if err != nil {
		t.Fatalf("Failed to create blank item: %v", err)
	}

	messageItemID := uuid.New()
	messageItem := &entities.MessageItem{
		MessageItemID: messageItemID,
		MeetingID:     meeting.MeetingID,
		ItemType:      "message",
		ItemOrder:     2,
		MeetingRole:   "Test Role",
		PrimaryText:   "Test Message",
	}
	err = service.CreateItem(testutil.TestCtx, messageItem)
	if err != nil {
		t.Fatalf("Failed to create message item: %v", err)
	}

	// Get all items
	items, err := service.GetItems(testutil.TestCtx)
	if err != nil {
		t.Fatalf("Failed to get items: %v", err)
	}
	if items == nil {
		t.Fatalf("Expected non-nil items slice")
	}
	if len(*items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(*items))
	}

	// Verify items are sorted by order
	if (*items)[0].GetOrder() != 1 {
		t.Errorf("Expected first item to have order 1, got %d", (*items)[0].GetOrder())
	}
	if (*items)[1].GetOrder() != 2 {
		t.Errorf("Expected second item to have order 2, got %d", (*items)[1].GetOrder())
	}

	// Find our test items in the list
	var foundBlankItem, foundMessageItem bool
	for _, item := range *items {
		switch v := item.(type) {
		case *entities.BlankItem:
			if v.BlankItemID == blankItemID {
				foundBlankItem = true
				if v.MeetingID != meeting.MeetingID {
					t.Errorf("Expected MeetingID %v, got %v", meeting.MeetingID, v.MeetingID)
				}
				if v.ItemType != "blank" {
					t.Errorf("Expected ItemType 'blank', got %v", v.ItemType)
				}
				if v.ItemOrder != 1 {
					t.Errorf("Expected ItemOrder 1, got %v", v.ItemOrder)
				}
				if v.MeetingRole != "Test Role" {
					t.Errorf("Expected MeetingRole 'Test Role', got %v", v.MeetingRole)
				}
			}
		case *entities.MessageItem:
			if v.MessageItemID == messageItemID {
				foundMessageItem = true
				if v.MeetingID != meeting.MeetingID {
					t.Errorf("Expected MeetingID %v, got %v", meeting.MeetingID, v.MeetingID)
				}
				if v.ItemType != "message" {
					t.Errorf("Expected ItemType 'message', got %v", v.ItemType)
				}
				if v.ItemOrder != 2 {
					t.Errorf("Expected ItemOrder 2, got %v", v.ItemOrder)
				}
				if v.MeetingRole != "Test Role" {
					t.Errorf("Expected MeetingRole 'Test Role', got %v", v.MeetingRole)
				}
				if v.PrimaryText != "Test Message" {
					t.Errorf("Expected PrimaryText 'Test Message', got %v", v.PrimaryText)
				}
			}
		}
	}
	if !foundBlankItem {
		t.Fatalf("Expected to find BlankItem with ID %v", blankItemID)
	}
	if !foundMessageItem {
		t.Fatalf("Expected to find MessageItem with ID %v", messageItemID)
	}
}

func TestCreateItemErrors(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)

	// Test creating an item with nil context
	err := service.CreateItem(nil, &entities.BlankItem{})
	if err == nil {
		t.Fatalf("Expected error when creating item with nil context")
	}

	// Test creating an item with invalid context (missing socialID)
	invalidCtx := context.Background()
	err = service.CreateItem(invalidCtx, &entities.BlankItem{})
	if err == nil {
		t.Fatalf("Expected error when creating item with invalid context")
	}

	// Test creating an item with nil item
	err = service.CreateItem(testutil.TestCtx, nil)
	if err == nil {
		t.Fatalf("Expected error when creating nil item")
	}

	// Test creating an item with invalid type
	type InvalidItem struct {
		entities.BlankItem
	}
	err = service.CreateItem(testutil.TestCtx, &InvalidItem{})
	if err == nil {
		t.Fatalf("Expected error when creating item with invalid type")
	}
}

func TestGetItemErrors(t *testing.T) {
	// Setup
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)
	ctx := context.Background()
	userID := uuid.New()
	meetingID := uuid.New()
	orgID := uuid.New()
	itemID := uuid.New()
	wrongOrgID := uuid.New()

	// Create test data
	createTestData(t, service, userID, meetingID, orgID, itemID)

	// Test cases
	tests := []struct {
		name          string
		ctx           context.Context
		itemID        uuid.UUID
		expectedError string
	}{
		{
			name:          "Invalid item type",
			ctx:           context.WithValue(ctx, helpers.SocialIDKey, "test-social-id"),
			itemID:        itemID,
			expectedError: "invalid item type",
		},
		{
			name:          "Wrong user",
			ctx:           context.WithValue(ctx, helpers.SocialIDKey, "wrong-social-id"),
			itemID:        itemID,
			expectedError: "item not found",
		},
		{
			name:          "Wrong organization",
			ctx:           context.WithValue(ctx, helpers.SocialIDKey, "test-social-id"),
			itemID:        itemID,
			expectedError: "item not found",
		},
		{
			name:          "Deleted item",
			ctx:           context.WithValue(ctx, helpers.SocialIDKey, "test-social-id"),
			itemID:        itemID,
			expectedError: "item not found",
		},
		{
			name:          "Deleted meeting",
			ctx:           context.WithValue(ctx, helpers.SocialIDKey, "test-social-id"),
			itemID:        itemID,
			expectedError: "item not found",
		},
		{
			name:          "Deleted organization",
			ctx:           context.WithValue(ctx, helpers.SocialIDKey, "test-social-id"),
			itemID:        itemID,
			expectedError: "item not found",
		},
		{
			name:          "Deleted user",
			ctx:           context.WithValue(ctx, helpers.SocialIDKey, "test-social-id"),
			itemID:        itemID,
			expectedError: "item not found",
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test case
			switch tt.name {
			case "Invalid item type":
				// Update the item to have an invalid type
				_, err := testutil.TestDB.Exec("UPDATE BlankItems SET item_type = 'invalid' WHERE id = ?", itemID)
				if err != nil {
					t.Fatalf("Failed to update item type: %v", err)
				}
				// NOTE: If this test fails, check GetItem implementation for invalid type handling.
			case "Wrong user":
				// Create a user in a different organization
				createUserInOrg(t, uuid.New(), uuid.New())
			case "Wrong organization":
				// Create a user in a different organization and update the item's meeting to be in that org
				createUserInOrg(t, uuid.New(), uuid.New())
				_, err := testutil.TestDB.Exec("UPDATE Meetings SET org_id = ? WHERE id = ?", wrongOrgID, meetingID)
				if err != nil {
					t.Fatalf("Failed to update meeting org: %v", err)
				}
			case "Deleted item":
				// Delete the item
				_, err := testutil.TestDB.Exec("UPDATE BlankItems SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ?", itemID)
				if err != nil {
					t.Fatalf("Failed to delete item: %v", err)
				}
			case "Deleted meeting":
				// Delete the meeting
				_, err := testutil.TestDB.Exec("UPDATE Meetings SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ?", meetingID)
				if err != nil {
					t.Fatalf("Failed to delete meeting: %v", err)
				}
			case "Deleted organization":
				// Delete the organization
				_, err := testutil.TestDB.Exec("UPDATE Organization SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ?", orgID)
				if err != nil {
					t.Fatalf("Failed to delete organization: %v", err)
				}
			case "Deleted user":
				// Delete the user
				_, err := testutil.TestDB.Exec("UPDATE Users SET deleted_dt = CURRENT_TIMESTAMP WHERE id = ?", userID)
				if err != nil {
					t.Fatalf("Failed to delete user: %v", err)
				}
			}

			// Call GetItem
			_, err := service.GetItem(tt.ctx, tt.itemID)

			// Check error
			if err == nil {
				t.Error("Expected error, got nil")
				return
			}
			if tt.name == "Deleted user" || tt.name == "Wrong user" {
				if !strings.Contains(err.Error(), tt.expectedError) && err.Error() != "sql: no rows in result set" {
					t.Errorf("Expected error containing %q or 'sql: no rows in result set', got %q", tt.expectedError, err.Error())
				}
				return
			}
			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing %q, got %q", tt.expectedError, err.Error())
			}
		})
	}
}

// Helper function to create test data
func createTestData(t *testing.T, service LowerThirdsService, userID, meetingID, orgID, itemID uuid.UUID) {
	// Create user
	user := &entities.User{
		UserID:   userID,
		Email:    "test@example.com",
		SocialID: null.StringFrom("test-social-id"),
	}
	err := service.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create organization
	org := &entities.Organization{
		OrgID: orgID,
		Name:  "Test Org " + orgID.String(),
	}
	ctx := context.WithValue(context.Background(), helpers.SocialIDKey, "test-social-id")
	err = service.CreateOrg(ctx, org)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}

	// Create org user
	_, err = testutil.TestDB.Exec(`
		INSERT INTO OrgUsers (user_id, org_id)
		VALUES (?, ?)
	`, userID, orgID)
	if err != nil {
		t.Fatalf("Failed to create org user: %v", err)
	}

	// Create meeting
	meeting := &entities.Meeting{
		MeetingID:   meetingID,
		OrgID:       orgID,
		Meeting:     "Test Meeting",
		MeetingDate: time.Now(),
	}
	err = service.CreateMeeting(ctx, meeting)
	if err != nil {
		t.Fatalf("Failed to create test meeting: %v", err)
	}

	// Create item
	blankItem := &entities.BlankItem{
		BlankItemID: itemID,
		MeetingID:   meetingID,
		ItemType:    "blank",
		ItemOrder:   1,
		MeetingRole: "Test Role",
	}
	err = service.CreateItem(ctx, blankItem)
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}
}

// Helper function to create a user in a different organization
func createUserInOrg(t *testing.T, userID, orgID uuid.UUID) {
	email := fmt.Sprintf("test+%s@example.com", userID.String())
	orgName := fmt.Sprintf("Test Org %s", orgID.String())
	_, err := testutil.TestDB.Exec(`
		INSERT INTO Users (id, email, first_name, last_name)
		VALUES (?, ?, 'Test', 'User')
	`, userID, email)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	_, err = testutil.TestDB.Exec(`
		INSERT INTO Organization (id, name)
		VALUES (?, ?)
	`, orgID, orgName)
	if err != nil {
		t.Fatalf("Failed to create organization: %v", err)
	}

	_, err = testutil.TestDB.Exec(`
		INSERT INTO OrgUsers (user_id, org_id)
		VALUES (?, ?)
	`, userID, orgID)
	if err != nil {
		t.Fatalf("Failed to create org user: %v", err)
	}
}

func TestSpeakerItemCRUD(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)

	_, _, meeting := testutil.CreateTestData(t, service)

	// Create
	speakerItemID := uuid.New()
	speakerItem := &entities.SpeakerItem{
		SpeakerItemID: speakerItemID,
		MeetingID:     meeting.MeetingID,
		ItemType:      "speaker",
		ItemOrder:     1,
		MeetingRole:   "Speaker Role",
		SpeakerName:   "Test Speaker",
	}
	err := service.CreateItem(testutil.TestCtx, speakerItem)
	if err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}

	// Get
	retrieved, err := service.GetItem(testutil.TestCtx, speakerItemID)
	if err != nil {
		t.Fatalf("GetItem failed: %v", err)
	}
	got, ok := retrieved.(*entities.SpeakerItem)
	if !ok {
		t.Fatalf("Expected SpeakerItem, got %T", retrieved)
	}
	if got.SpeakerItemID != speakerItemID || got.MeetingID != meeting.MeetingID || got.ItemType != "speaker" || got.ItemOrder != 1 || got.MeetingRole != "Speaker Role" || got.SpeakerName != "Test Speaker" {
		t.Errorf("SpeakerItem fields do not match expected values")
	}

	// Update
	got.SpeakerName = "Updated Speaker"
	err = service.UpdateItem(testutil.TestCtx, speakerItemID, got)
	if err != nil {
		t.Fatalf("UpdateItem failed: %v", err)
	}
	updated, err := service.GetItem(testutil.TestCtx, speakerItemID)
	if err != nil {
		t.Fatalf("GetItem after update failed: %v", err)
	}
	updatedSpeaker, ok := updated.(*entities.SpeakerItem)
	if !ok {
		t.Fatalf("Expected SpeakerItem after update, got %T", updated)
	}
	if updatedSpeaker.SpeakerName != "Updated Speaker" {
		t.Errorf("Expected SpeakerName 'Updated Speaker', got %v", updatedSpeaker.SpeakerName)
	}

	// Delete
	err = service.DeleteItem(testutil.TestCtx, speakerItemID)
	if err != nil {
		t.Fatalf("DeleteItem failed: %v", err)
	}
	_, err = service.GetItem(testutil.TestCtx, speakerItemID)
	if err == nil {
		t.Fatalf("Expected error when getting deleted SpeakerItem, got nil")
	}
}

func TestLyricsItemCRUD(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)

	_, _, meeting := testutil.CreateTestData(t, service)

	// Create
	lyricsItemID := uuid.New()
	lyricsItem := &entities.LyricsItem{
		LyricsItemID:    lyricsItemID,
		MeetingID:       meeting.MeetingID,
		ItemType:        "lyrics",
		ItemOrder:       1,
		MeetingRole:     "Lyrics Role",
		HymnID:          "H123",
		ShowTranslation: true,
	}
	err := service.CreateItem(testutil.TestCtx, lyricsItem)
	if err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}

	// Get
	retrieved, err := service.GetItem(testutil.TestCtx, lyricsItemID)
	if err != nil {
		t.Fatalf("GetItem failed: %v", err)
	}
	got, ok := retrieved.(*entities.LyricsItem)
	if !ok {
		t.Fatalf("Expected LyricsItem, got %T", retrieved)
	}
	if got.LyricsItemID != lyricsItemID || got.MeetingID != meeting.MeetingID || got.ItemType != "lyrics" || got.ItemOrder != 1 || got.MeetingRole != "Lyrics Role" || got.HymnID != "H123" || !got.ShowTranslation {
		t.Errorf("LyricsItem fields do not match expected values")
	}

	// Update
	got.HymnID = "H456"
	got.ShowTranslation = false
	err = service.UpdateItem(testutil.TestCtx, lyricsItemID, got)
	if err != nil {
		t.Fatalf("UpdateItem failed: %v", err)
	}
	updated, err := service.GetItem(testutil.TestCtx, lyricsItemID)
	if err != nil {
		t.Fatalf("GetItem after update failed: %v", err)
	}
	updatedLyrics, ok := updated.(*entities.LyricsItem)
	if !ok {
		t.Fatalf("Expected LyricsItem after update, got %T", updated)
	}
	if updatedLyrics.HymnID != "H456" || updatedLyrics.ShowTranslation != false {
		t.Errorf("Expected HymnID 'H456' and ShowTranslation false, got %v and %v", updatedLyrics.HymnID, updatedLyrics.ShowTranslation)
	}

	// Delete
	err = service.DeleteItem(testutil.TestCtx, lyricsItemID)
	if err != nil {
		t.Fatalf("DeleteItem failed: %v", err)
	}
	_, err = service.GetItem(testutil.TestCtx, lyricsItemID)
	if err == nil {
		t.Fatalf("Expected error when getting deleted LyricsItem, got nil")
	}
}

func TestTimerItemCRUD(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)

	_, _, meeting := testutil.CreateTestData(t, service)

	// Create
	timerItemID := uuid.New()
	timerItem := &entities.TimerItem{
		TimerItemID:        timerItemID,
		MeetingID:          meeting.MeetingID,
		ItemType:           "timer",
		ItemOrder:          1,
		MeetingRole:        "Timer Role",
		ShowMeetingDetails: true,
	}
	err := service.CreateItem(testutil.TestCtx, timerItem)
	if err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}

	// Get
	retrieved, err := service.GetItem(testutil.TestCtx, timerItemID)
	if err != nil {
		t.Fatalf("GetItem failed: %v", err)
	}
	got, ok := retrieved.(*entities.TimerItem)
	if !ok {
		t.Fatalf("Expected TimerItem, got %T", retrieved)
	}
	if got.TimerItemID != timerItemID || got.MeetingID != meeting.MeetingID || got.ItemType != "timer" || got.ItemOrder != 1 || got.MeetingRole != "Timer Role" || !got.ShowMeetingDetails {
		t.Errorf("TimerItem fields do not match expected values")
	}

	// Update
	got.ShowMeetingDetails = false
	err = service.UpdateItem(testutil.TestCtx, timerItemID, got)
	if err != nil {
		t.Fatalf("UpdateItem failed: %v", err)
	}
	updated, err := service.GetItem(testutil.TestCtx, timerItemID)
	if err != nil {
		t.Fatalf("GetItem after update failed: %v", err)
	}
	updatedTimer, ok := updated.(*entities.TimerItem)
	if !ok {
		t.Fatalf("Expected TimerItem after update, got %T", updated)
	}
	if updatedTimer.ShowMeetingDetails != false {
		t.Errorf("Expected ShowMeetingDetails false, got %v", updatedTimer.ShowMeetingDetails)
	}

	// Delete
	err = service.DeleteItem(testutil.TestCtx, timerItemID)
	if err != nil {
		t.Fatalf("DeleteItem failed: %v", err)
	}
	_, err = service.GetItem(testutil.TestCtx, timerItemID)
	if err == nil {
		t.Fatalf("Expected error when getting deleted TimerItem, got nil")
	}
}

func TestUpdateNonExistentItem(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)

	nonExistentID := uuid.New()
	blankItem := &entities.BlankItem{
		BlankItemID: nonExistentID,
		ItemType:    "blank",
		ItemOrder:   1,
		MeetingRole: "Test Role",
	}
	err := service.UpdateItem(testutil.TestCtx, nonExistentID, blankItem)
	if err == nil {
		t.Fatalf("Expected error when updating non-existent item, got nil")
	}
	if err.Error() != "sql: no rows in result set" {
		t.Fatalf("Expected 'sql: no rows in result set' error, got: %v", err)
	}
}

func TestDeleteNonExistentItem(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)

	nonExistentID := uuid.New()
	err := service.DeleteItem(testutil.TestCtx, nonExistentID)
	if err == nil {
		t.Fatalf("Expected error when deleting non-existent item, got nil")
	}
}

func TestCreateItemWithDuplicateID(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)

	_, _, meeting := testutil.CreateTestData(t, service)

	itemID := uuid.New()
	blankItem := &entities.BlankItem{
		BlankItemID: itemID,
		MeetingID:   meeting.MeetingID,
		ItemType:    "blank",
		ItemOrder:   1,
		MeetingRole: "Test Role",
	}
	err := service.CreateItem(testutil.TestCtx, blankItem)
	if err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}

	// Try to create another item with the same ID
	duplicateItem := &entities.BlankItem{
		BlankItemID: itemID,
		MeetingID:   meeting.MeetingID,
		ItemType:    "blank",
		ItemOrder:   2,
		MeetingRole: "Another Role",
	}
	err = service.CreateItem(testutil.TestCtx, duplicateItem)
	if err == nil {
		t.Fatalf("Expected error when creating item with duplicate ID, got nil")
	}
}

func TestGetItemsByMeetingWithNoItems(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)

	_, _, meeting := testutil.CreateTestData(t, service)

	items, err := service.GetItemsByMeeting(testutil.TestCtx, meeting.MeetingID)
	if err != nil {
		t.Fatalf("GetItemsByMeeting failed: %v", err)
	}
	if items == nil {
		t.Fatalf("Expected non-nil items slice")
	}
	if len(*items) != 0 {
		t.Fatalf("Expected 0 items, got %d", len(*items))
	}
}

func TestGetItemsErrors(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)

	// Test getting items with nil context
	_, err := service.GetItems(nil)
	if err == nil {
		t.Fatalf("Expected error when getting items with nil context")
	}

	// Test getting items with invalid context (missing socialID)
	invalidCtx := context.Background()
	_, err = service.GetItems(invalidCtx)
	if err == nil {
		t.Fatalf("Expected error when getting items with invalid context")
	}

	// Test getting items with non-existent user
	ctxWrongUser := context.WithValue(testutil.TestCtx, helpers.SocialIDKey, "non-existent-social-id")
	_, err = service.GetItems(ctxWrongUser)
	if err == nil {
		t.Fatalf("Expected error when getting items with non-existent user")
	}
}

func TestGetItemsByMeetingErrors(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()
	service := New(testutil.TestDB, testutil.TestLogger)

	// Test getting items with nil context
	_, err := service.GetItemsByMeeting(nil, uuid.New())
	if err == nil {
		t.Fatalf("Expected error when getting items with nil context")
	}

	// Test getting items with invalid context (missing socialID)
	invalidCtx := context.Background()
	_, err = service.GetItemsByMeeting(invalidCtx, uuid.New())
	if err == nil {
		t.Fatalf("Expected error when getting items with invalid context")
	}

	// Test getting items with non-existent user
	ctxWrongUser := context.WithValue(testutil.TestCtx, helpers.SocialIDKey, "non-existent-social-id")
	_, err = service.GetItemsByMeeting(ctxWrongUser, uuid.New())
	if err == nil {
		t.Fatalf("Expected error when getting items with non-existent user")
	}

	// Test getting items for non-existent meeting
	_, err = service.GetItemsByMeeting(testutil.TestCtx, uuid.New())
	if err == nil {
		t.Fatalf("Expected error when getting items for non-existent meeting")
	}
	if err.Error() != "sql: no rows in result set" {
		t.Fatalf("Expected 'sql: no rows in result set' error, got: %v", err)
	}
}
