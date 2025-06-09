package storage

import (
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/storage/testutil"
	"testing"

	"github.com/google/uuid"
)

func TestCreateAndGetItem(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()

	service := New(testutil.TestDB, testutil.TestLogger)

	// Setup test data
	_, _, meeting := testutil.SetupTestData(t)

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

	// Test CreateItem
	err := service.CreateItem(testutil.TestCtx, messageItem)
	if err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}

	// Test GetItem
	retrievedItem, err := service.GetItem(testutil.TestCtx, messageItemID)
	if err != nil {
		t.Fatalf("GetItem failed: %v", err)
	}

	// Verify item data
	retrievedMessageItem, ok := retrievedItem.(*entities.MessageItem)
	if !ok {
		t.Fatalf("Expected MessageItem, got %T", retrievedItem)
	}

	if retrievedMessageItem.MessageItemID != messageItem.MessageItemID {
		t.Errorf("Expected MessageItemID %v, got %v", messageItem.MessageItemID, retrievedMessageItem.MessageItemID)
	}
	if retrievedMessageItem.MeetingID != messageItem.MeetingID {
		t.Errorf("Expected MeetingID %v, got %v", messageItem.MeetingID, retrievedMessageItem.MeetingID)
	}
	if retrievedMessageItem.ItemType != messageItem.ItemType {
		t.Errorf("Expected ItemType %v, got %v", messageItem.ItemType, retrievedMessageItem.ItemType)
	}
	if retrievedMessageItem.ItemOrder != messageItem.ItemOrder {
		t.Errorf("Expected ItemOrder %v, got %v", messageItem.ItemOrder, retrievedMessageItem.ItemOrder)
	}
	if retrievedMessageItem.MeetingRole != messageItem.MeetingRole {
		t.Errorf("Expected MeetingRole %v, got %v", messageItem.MeetingRole, retrievedMessageItem.MeetingRole)
	}
	if retrievedMessageItem.PrimaryText != messageItem.PrimaryText {
		t.Errorf("Expected PrimaryText %v, got %v", messageItem.PrimaryText, retrievedMessageItem.PrimaryText)
	}
}

func TestDeleteItem(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()

	service := New(testutil.TestDB, testutil.TestLogger)

	// Setup test data
	_, _, meeting := testutil.SetupTestData(t)

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

	// Test DeleteItem
	err = service.DeleteItem(testutil.TestCtx, messageItemID)
	if err != nil {
		t.Fatalf("DeleteItem failed: %v", err)
	}

	// Verify item is deleted
	retrievedItem, err := service.GetItem(testutil.TestCtx, messageItemID)
	if err != nil {
		t.Fatalf("GetItem failed: %v", err)
	}
	if retrievedItem != nil {
		t.Errorf("Expected nil item after deletion, got %v", retrievedItem)
	}
}

func TestUpdateItem(t *testing.T) {
	testutil.SetupTest(t)
	defer testutil.TeardownTest()

	service := New(testutil.TestDB, testutil.TestLogger)

	// Setup test data
	_, _, meeting := testutil.SetupTestData(t)

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

	// Update the item
	updatedMessageItem := &entities.MessageItem{
		MessageItemID: messageItemID,
		MeetingID:     meeting.MeetingID,
		ItemType:      "message",
		ItemOrder:     2,
		MeetingRole:   "Updated Role",
		PrimaryText:   "Updated Message",
	}

	// Test UpdateItem
	err = service.UpdateItem(testutil.TestCtx, messageItemID, updatedMessageItem)
	if err != nil {
		t.Fatalf("UpdateItem failed: %v", err)
	}

	// Verify item is updated
	retrievedItem, err := service.GetItem(testutil.TestCtx, messageItemID)
	if err != nil {
		t.Fatalf("GetItem failed: %v", err)
	}

	retrievedMessageItem, ok := retrievedItem.(*entities.MessageItem)
	if !ok {
		t.Fatalf("Expected MessageItem, got %T", retrievedItem)
	}

	if retrievedMessageItem.ItemOrder != updatedMessageItem.ItemOrder {
		t.Errorf("Expected ItemOrder %v, got %v", updatedMessageItem.ItemOrder, retrievedMessageItem.ItemOrder)
	}
	if retrievedMessageItem.MeetingRole != updatedMessageItem.MeetingRole {
		t.Errorf("Expected MeetingRole %v, got %v", updatedMessageItem.MeetingRole, retrievedMessageItem.MeetingRole)
	}
	if retrievedMessageItem.PrimaryText != updatedMessageItem.PrimaryText {
		t.Errorf("Expected PrimaryText %v, got %v", updatedMessageItem.PrimaryText, retrievedMessageItem.PrimaryText)
	}
}
