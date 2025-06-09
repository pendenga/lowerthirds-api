package entities

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestBlankItemUnmarshal(t *testing.T) {
	blankJSON := `{
		"id": "11111111-1111-1111-1111-111111111111",
		"meeting_id": "22222222-2222-2222-2222-222222222222",
		"type": "blank",
		"order": 1,
		"meeting_role": "my role"
	}`
	var blank BlankItem
	if err := json.Unmarshal([]byte(blankJSON), &blank); err != nil {
		t.Errorf("failed to unmarshal BlankItem: %v", err)
	}
	if blank.BlankItemID != uuid.MustParse("11111111-1111-1111-1111-111111111111") {
		t.Errorf("unexpected BlankItemID: %v", blank.BlankItemID)
	}
	if blank.MeetingID != uuid.MustParse("22222222-2222-2222-2222-222222222222") {
		t.Errorf("unexpected BlankItemID: %v", blank.MeetingID)
	}
}

func TestBlankItemUnmarshal_new(t *testing.T) {
	blankJSON := `{
 		"meeting_id": "22222222-2222-2222-2222-222222222222",
		"type": "blank",
		"order": 1,
		"meeting_role": "my role"
	}`
	var blank BlankItem
	if err := json.Unmarshal([]byte(blankJSON), &blank); err != nil {
		t.Errorf("failed to unmarshal BlankItem: %v", err)
	}
	if blank.BlankItemID != uuid.Nil {
		t.Errorf("expected BlankItemID to be zero UUID, got: %v", blank.BlankItemID)
	}
	if blank.MeetingID != uuid.MustParse("22222222-2222-2222-2222-222222222222") {
		t.Errorf("unexpected BlankItemID: %v", blank.MeetingID)
	}
}

func TestMessageItemUnmarshal(t *testing.T) {
	messageJSON := `{
		"id": "33333333-3333-3333-3333-333333333333",
		"meeting_id": "44444444-4444-4444-4444-444444444444",
		"type": "message",
		"order": 2,
		"meeting_role": "my role",
		"primary_text": "Hello"
	}`
	var message MessageItem
	if err := json.Unmarshal([]byte(messageJSON), &message); err != nil {
		t.Errorf("failed to unmarshal MessageItem: %v", err)
	}
	if message.PrimaryText != "Hello" {
		t.Errorf("unexpected PrimaryText: %v", message.PrimaryText)
	}
	if message.MessageItemID != uuid.MustParse("33333333-3333-3333-3333-333333333333") {
		t.Errorf("unexpected BlankItemID: %v", message.MessageItemID)
	}
	if message.MeetingID != uuid.MustParse("44444444-4444-4444-4444-444444444444") {
		t.Errorf("unexpected BlankItemID: %v", message.MeetingID)
	}
}

func TestMessageItemUnmarshal_new(t *testing.T) {
	messageJSON := `{
		"meeting_id": "44444444-4444-4444-4444-444444444444",
		"type": "message",
		"order": 2,
		"meeting_role": "my role",
		"primary_text": "Hello"
	}`
	var message MessageItem
	if err := json.Unmarshal([]byte(messageJSON), &message); err != nil {
		t.Errorf("failed to unmarshal MessageItem: %v", err)
	}
	if message.PrimaryText != "Hello" {
		t.Errorf("unexpected PrimaryText: %v", message.PrimaryText)
	}
	if message.MessageItemID != uuid.Nil {
		t.Errorf("expected BlankItemID to be zero UUID, got: %v", message.MessageItemID)
	}
	if message.MeetingID != uuid.MustParse("44444444-4444-4444-4444-444444444444") {
		t.Errorf("unexpected BlankItemID: %v", message.MeetingID)
	}
}

func TestParseItemJSON_BlankItem(t *testing.T) {
	blankJSON := `{
		"id": "11111111-1111-1111-1111-111111111111",
		"meeting_id": "22222222-2222-2222-2222-222222222222",
		"type": "blank",
		"order": 1,
		"meeting_role": "my role"
	}`
	item, err := ParseItemJSON([]byte(blankJSON))
	if err != nil {
		t.Fatalf("ParseItemJSON failed: %v", err)
	}
	blank, ok := item.(*BlankItem)
	if !ok {
		t.Fatalf("expected *BlankItem, got %T", item)
	}
	if blank.BlankItemID != uuid.MustParse("11111111-1111-1111-1111-111111111111") {
		t.Errorf("unexpected BlankItemID: %v", blank.BlankItemID)
	}
	if blank.MeetingID != uuid.MustParse("22222222-2222-2222-2222-222222222222") {
		t.Errorf("unexpected MeetingID: %v", blank.MeetingID)
	}
	if blank.ItemType != "blank" {
		t.Errorf("unexpected ItemType: %v", blank.ItemType)
	}
	if blank.ItemOrder != 1 {
		t.Errorf("unexpected ItemOrder: %v", blank.ItemOrder)
	}
	if blank.MeetingRole != "my role" {
		t.Errorf("unexpected MeetingRole: %v", blank.MeetingRole)
	}
}
