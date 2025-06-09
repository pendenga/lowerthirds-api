package storage

import (
	"context"
	"database/sql"
	"errors"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
	"sort"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

func (s lowerThirdsService) CreateItem(ctx context.Context, item entities.Item) error {
	if ctx == nil {
		return errors.New("context is required")
	}
	if item == nil {
		return errors.New("item is required")
	}

	socialID, ok := ctx.Value(helpers.SocialIDKey).(string)
	if !ok {
		return errors.New("socialID is required in context")
	}
	s.logger.Debug("[CreateItem] for socialID ", socialID)
	s.logger.Debugf("[CreateItem] %+v", item)

	// Each type of item handled separately
	switch v := item.(type) {
	case *entities.BlankItem:
		if v.BlankItemID == uuid.Nil {
			v.BlankItemID = uuid.New()
		}
		s.logger.Debugf("[CreateItem] createBlankItem %+v", v)
		err := s.createBlankItem(v)
		if err != nil {
			s.logger.Error("error creating blankItem ", err)
			return err
		}
	case *entities.LyricsItem:
		if v.LyricsItemID == uuid.Nil {
			v.LyricsItemID = uuid.New()
		}
		s.logger.Debugf("[CreateItem] createLyricsItem %+v", v)
		err := s.createLyricsItem(v)
		if err != nil {
			s.logger.Error("error creating lyricsItem ", err)
			return err
		}
	case *entities.MessageItem:
		if v.MessageItemID == uuid.Nil {
			v.MessageItemID = uuid.New()
		}
		s.logger.Debugf("[CreateItem] createMessageItem %+v", v)
		err := s.createMessageItem(v)
		if err != nil {
			s.logger.Error("error creating messageItem ", err)
			return err
		}
	case *entities.SpeakerItem:
		if v.SpeakerItemID == uuid.Nil {
			v.SpeakerItemID = uuid.New()
		}
		s.logger.Debugf("[CreateItem] createSpeakerItem %+v", v)
		err := s.createSpeakerItem(v)
		if err != nil {
			s.logger.Error("error creating speakerItem ", err)
			return err
		}
	case *entities.TimerItem:
		if v.TimerItemID == uuid.Nil {
			v.TimerItemID = uuid.New()
		}
		s.logger.Debugf("[CreateItem] createTimerItem %+v", v)
		err := s.createTimerItem(v)
		if err != nil {
			s.logger.Error("error creating timerItem ", err)
			return err
		}
	default:
		return errors.New("unsupported item type")
	}
	return nil
}

func (s lowerThirdsService) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	socialID := ctx.Value(helpers.SocialIDKey).(string)
	s.logger.Debug("DeleteItems for socialID ", socialID, " itemID ", itemID)
	user, err := s.GetUserBySocialID(ctx, socialID)
	if err != nil {
		s.logger.Error("User not found by socialID", err)
		return err
	}
	s.logger.Debug("DeleteItems for userID ", user.UserID, " itemID ", itemID)

	var totalAffectedRows int64 = 0

	// Query each type of item separately
	affectedRows, err := s.deleteBlankItem(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error deleting blankItem ", err)
		return err
	}
	totalAffectedRows = totalAffectedRows + affectedRows
	affectedRows, err = s.deleteMessageItem(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error deleting messageItem ", err)
		return err
	}
	totalAffectedRows = totalAffectedRows + affectedRows
	affectedRows, err = s.deleteSpeakerItem(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error deleting speakerItem ", err)
		return err
	}
	totalAffectedRows = totalAffectedRows + affectedRows
	affectedRows, err = s.deleteLyricsItem(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error deleting lyricsItem ", err)
		return err
	}
	totalAffectedRows = totalAffectedRows + affectedRows
	affectedRows, err = s.deleteTimerItem(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error deleting timerItem ", err)
		return err
	}
	totalAffectedRows = totalAffectedRows + affectedRows

	s.logger.Info("DeleteItems affectedRows rows: ", totalAffectedRows)

	return nil
}

func (s lowerThirdsService) GetItem(ctx context.Context, itemID uuid.UUID) (entities.Item, error) {
	socialID := ctx.Value(helpers.SocialIDKey).(string)
	s.logger.Debug("GetItem for socialID ", socialID, " itemID ", itemID)
	user, err := s.GetUserBySocialID(ctx, socialID)
	if err != nil {
		s.logger.Error("User not found by socialID", err)
		return nil, err
	}
	s.logger.Debug("GetItem for userID ", user.UserID, " itemID ", itemID)

	// Query each type of item separately
	blankItem, err := s.getBlankItemByID(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error querying blankItem ", err)
		return nil, err
	}
	if blankItem != nil {
		if blankItem.ItemType != "blank" {
			return nil, errors.New("invalid item type")
		}
		return blankItem, nil
	}
	messageItem, err := s.getMessageItemByID(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error querying messageItem ", err)
		return nil, err
	}
	if messageItem != nil {
		if messageItem.ItemType != "message" {
			return nil, errors.New("invalid item type")
		}
		return messageItem, nil
	}
	speakerItem, err := s.getSpeakerItemByID(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error querying speakerItem ", err)
		return nil, err
	}
	if speakerItem != nil {
		if speakerItem.ItemType != "speaker" {
			return nil, errors.New("invalid item type")
		}
		return speakerItem, nil
	}
	lyricsItem, err := s.getLyricsItemByID(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error querying lyricsItem ", err)
		return nil, err
	}
	if lyricsItem != nil {
		if lyricsItem.ItemType != "lyrics" {
			return nil, errors.New("invalid item type")
		}
		return lyricsItem, nil
	}
	timerItem, err := s.getTimerItemByID(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error querying timerItem ", err)
		return nil, err
	}
	if timerItem != nil {
		if timerItem.ItemType != "timer" {
			return nil, errors.New("invalid item type")
		}
		return timerItem, nil
	}
	return nil, errors.New("item not found")
}

func (s lowerThirdsService) getAllItemsByUser(userID uuid.UUID) ([]entities.Item, error) {
	s.logger.Debug("getAllItemsByUser for userID ", userID)

	query := `
		WITH user_meetings AS (
			SELECT m.id as meeting_id
			FROM OrgUsers ou
			INNER JOIN Users u ON u.id = ou.user_id AND u.deleted_dt IS NULL
			INNER JOIN Organization o ON o.id = ou.org_id AND o.deleted_dt IS NULL
			INNER JOIN Meetings m ON m.org_id = ou.org_id AND m.deleted_dt IS NULL
			WHERE ou.user_id = ? AND ou.deleted_dt IS NULL
		)
		SELECT 
			id, meeting_id, meeting_role, item_type, item_order,
			NULL as primary_text, NULL as secondary_text,
			NULL as speaker_name, NULL as title, NULL as expected_duration,
			NULL as hymn_id, NULL as show_translation,
			NULL as show_meeting_details,
			'blank' as source_table
		FROM BlankItems
		WHERE meeting_id IN (SELECT meeting_id FROM user_meetings)
		AND deleted_dt IS NULL
		UNION ALL
		SELECT 
			id, meeting_id, meeting_role, item_type, item_order,
			primary_text, secondary_text,
			NULL as speaker_name, NULL as title, NULL as expected_duration,
			NULL as hymn_id, NULL as show_translation,
			NULL as show_meeting_details,
			'message' as source_table
		FROM MessageItems
		WHERE meeting_id IN (SELECT meeting_id FROM user_meetings)
		AND deleted_dt IS NULL
		UNION ALL
		SELECT 
			id, meeting_id, meeting_role, item_type, item_order,
			NULL as primary_text, NULL as secondary_text,
			speaker_name, title, expected_duration,
			NULL as hymn_id, NULL as show_translation,
			NULL as show_meeting_details,
			'speaker' as source_table
		FROM SpeakerItems
		WHERE meeting_id IN (SELECT meeting_id FROM user_meetings)
		AND deleted_dt IS NULL
		UNION ALL
		SELECT 
			id, meeting_id, meeting_role, item_type, item_order,
			NULL as primary_text, NULL as secondary_text,
			NULL as speaker_name, NULL as title, NULL as expected_duration,
			hymn_id, show_translation,
			NULL as show_meeting_details,
			'lyrics' as source_table
		FROM LyricsItems
		WHERE meeting_id IN (SELECT meeting_id FROM user_meetings)
		AND deleted_dt IS NULL
		UNION ALL
		SELECT 
			id, meeting_id, meeting_role, item_type, item_order,
			NULL as primary_text, NULL as secondary_text,
			NULL as speaker_name, NULL as title, NULL as expected_duration,
			NULL as hymn_id, NULL as show_translation,
			show_meeting_details,
			'timer' as source_table
		FROM TimerItems
		WHERE meeting_id IN (SELECT meeting_id FROM user_meetings)
		AND deleted_dt IS NULL
		ORDER BY item_order`

	rows, err := s.MySqlDB.Query(query, userID)
	if err != nil {
		s.logger.Error("Error querying all items: ", err)
		return nil, err
	}
	defer rows.Close()

	var items []entities.Item
	for rows.Next() {
		var (
			id                 uuid.UUID
			meetingID          uuid.UUID
			meetingRole        string
			itemType           string
			itemOrder          int
			primaryText        sql.NullString
			secondaryText      sql.NullString
			speakerName        sql.NullString
			title              sql.NullString
			expectedDuration   sql.NullInt32
			hymnID             sql.NullString
			showTranslation    sql.NullBool
			showMeetingDetails sql.NullBool
			sourceTable        string
		)

		err := rows.Scan(
			&id, &meetingID, &meetingRole, &itemType, &itemOrder,
			&primaryText, &secondaryText,
			&speakerName, &title, &expectedDuration,
			&hymnID, &showTranslation,
			&showMeetingDetails,
			&sourceTable,
		)
		if err != nil {
			s.logger.Error("Error scanning row: ", err)
			return nil, err
		}

		switch sourceTable {
		case "blank":
			items = append(items, &entities.BlankItem{
				BlankItemID: id,
				MeetingID:   meetingID,
				ItemType:    itemType,
				ItemOrder:   itemOrder,
				MeetingRole: meetingRole,
			})
		case "message":
			items = append(items, &entities.MessageItem{
				MessageItemID: id,
				MeetingID:     meetingID,
				ItemType:      itemType,
				ItemOrder:     itemOrder,
				MeetingRole:   meetingRole,
				PrimaryText:   primaryText.String,
				SecondaryText: null.StringFromPtr(&secondaryText.String),
			})
		case "speaker":
			var expectedDurationPtr *int64
			if expectedDuration.Valid {
				duration := int64(expectedDuration.Int32)
				expectedDurationPtr = &duration
			}
			items = append(items, &entities.SpeakerItem{
				SpeakerItemID:    id,
				MeetingID:        meetingID,
				ItemType:         itemType,
				ItemOrder:        itemOrder,
				MeetingRole:      meetingRole,
				SpeakerName:      speakerName.String,
				Title:            null.StringFromPtr(&title.String),
				ExpectedDuration: null.IntFromPtr(expectedDurationPtr),
			})
		case "lyrics":
			items = append(items, &entities.LyricsItem{
				LyricsItemID:    id,
				MeetingID:       meetingID,
				ItemType:        itemType,
				ItemOrder:       itemOrder,
				MeetingRole:     meetingRole,
				HymnID:          hymnID.String,
				ShowTranslation: showTranslation.Bool,
			})
		case "timer":
			items = append(items, &entities.TimerItem{
				TimerItemID:        id,
				MeetingID:          meetingID,
				ItemType:           itemType,
				ItemOrder:          itemOrder,
				MeetingRole:        meetingRole,
				ShowMeetingDetails: showMeetingDetails.Bool,
			})
		}
	}

	if err = rows.Err(); err != nil {
		s.logger.Error("Error iterating rows: ", err)
		return nil, err
	}

	return items, nil
}

func (s lowerThirdsService) GetItems(ctx context.Context) (*[]entities.Item, error) {
	if ctx == nil {
		return nil, errors.New("context is required")
	}

	socialID, ok := ctx.Value(helpers.SocialIDKey).(string)
	if !ok {
		return nil, errors.New("socialID is required in context")
	}
	s.logger.Debug("GetItems for socialID ", socialID)
	user, err := s.GetUserBySocialID(ctx, socialID)
	if err != nil {
		s.logger.Error("User not found by socialID", err)
		return nil, err
	}
	s.logger.Debug("GetItems for userID ", user.UserID)

	items, err := s.getAllItemsByUser(user.UserID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &items, nil
}

func (s lowerThirdsService) GetItemsByMeeting(ctx context.Context, meetingID uuid.UUID) (*[]entities.Item, error) {
	if ctx == nil {
		return nil, errors.New("context is required")
	}

	socialID, ok := ctx.Value(helpers.SocialIDKey).(string)
	if !ok {
		return nil, errors.New("socialID is required in context")
	}
	s.logger.Debug("GetItem for socialID ", socialID, " meetingID ", meetingID)
	user, err := s.GetUserBySocialID(ctx, socialID)
	if err != nil {
		s.logger.Error("User not found by socialID", err)
		return nil, err
	}
	s.logger.Debug("GetItem for userID ", user.UserID, " meetingID ", meetingID)

	// Query each type of item separately
	blankItems, err := s.getBlankItemsByMeeting(user.UserID, meetingID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	messageItems, err := s.getMessageItemsByMeeting(user.UserID, meetingID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	speakerItems, err := s.getSpeakerItemsByMeeting(user.UserID, meetingID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	lyricsItems, err := s.getLyricsItemsByMeeting(user.UserID, meetingID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	timerItems, err := s.getTimerItemsByMeeting(user.UserID, meetingID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	// Combine in an Item interface set
	var allItems []entities.Item
	for _, b := range blankItems {
		allItems = append(allItems, &b)
	}
	for _, m := range messageItems {
		allItems = append(allItems, &m)
	}
	for _, s := range speakerItems {
		allItems = append(allItems, &s)
	}
	for _, l := range lyricsItems {
		allItems = append(allItems, &l)
	}
	for _, t := range timerItems {
		allItems = append(allItems, &t)
	}

	// Sort items by GetOrder
	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].GetOrder() < allItems[j].GetOrder()
	})

	return &allItems, nil
}

func (s lowerThirdsService) UpdateItem(ctx context.Context, itemID uuid.UUID, item entities.Item) error {
	socialID := ctx.Value(helpers.SocialIDKey).(string)
	s.logger.Debug("UpdateItems for socialID ", socialID, " itemID ", itemID)

	// TODO: somehow check that the ID being used matches up with an existing record of the same type.
	// Theoretically we could get duplicate ID numbers in different tables

	// Each type of item handled separately
	switch v := item.(type) {
	case *entities.BlankItem:
		if v.BlankItemID == uuid.Nil {
			v.BlankItemID = itemID
		}
		err := s.updateBlankItem(itemID, v)
		if err != nil {
			s.logger.Error("error updating blankItem ", err)
			return err
		}
		return nil
	case *entities.LyricsItem:
		if v.LyricsItemID == uuid.Nil {
			v.LyricsItemID = itemID
		}
		err := s.updateLyricsItem(itemID, v)
		if err != nil {
			s.logger.Error("error updating lyricsItem ", err)
			return err
		}
		return nil
	case *entities.MessageItem:
		if v.MessageItemID == uuid.Nil {
			v.MessageItemID = itemID
		}
		err := s.updateMessageItem(itemID, v)
		if err != nil {
			s.logger.Error("error updating messageItem ", err)
			return err
		}
		return nil
	case *entities.SpeakerItem:
		if v.SpeakerItemID == uuid.Nil {
			v.SpeakerItemID = itemID
		}
		err := s.updateSpeakerItem(itemID, v)
		if err != nil {
			s.logger.Error("error updating speakerItem ", err)
			return err
		}
		return nil
	case *entities.TimerItem:
		if v.TimerItemID == uuid.Nil {
			v.TimerItemID = itemID
		}
		err := s.updateTimerItem(itemID, v)
		if err != nil {
			s.logger.Error("error updating timerItem ", err)
			return err
		}
		return nil
	}
	return errors.New("unsupported item type")
}
