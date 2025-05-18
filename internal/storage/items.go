package storage

import (
	"context"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
	"sort"
)

func (s lowerThirdsService) CreateItem(ctx context.Context, item entities.Item) error {
	socialID := ctx.Value(helpers.SocialIDKey).(string)
	s.logger.Debug("CreateItems for socialID ", socialID)

	// Each type of item handled separately
	switch v := item.(type) {
	case *entities.BlankItem:
		if v.BlankItemID == uuid.Nil {
			v.BlankItemID = uuid.New()
		}
		err := s.createBlankItem(v)
		if err != nil {
			s.logger.Error("error creating blankItem ", err)
			return err
		}
		break
	case *entities.LyricsItem:
		if v.LyricsItemID == uuid.Nil {
			v.LyricsItemID = uuid.New()
		}
		err := s.createLyricsItem(v)
		if err != nil {
			s.logger.Error("error creating lyricsItem ", err)
			return err
		}
		break
	case *entities.MessageItem:
		if v.MessageItemID == uuid.Nil {
			v.MessageItemID = uuid.New()
		}
		err := s.createMessageItem(v)
		if err != nil {
			s.logger.Error("error creating messageItem ", err)
			return err
		}
		break
	case *entities.SpeakerItem:
		if v.SpeakerItemID == uuid.Nil {
			v.SpeakerItemID = uuid.New()
		}
		err := s.createSpeakerItem(v)
		if err != nil {
			s.logger.Error("error creating speakerItem ", err)
			return err
		}
		break
	case *entities.TimerItem:
		if v.TimerItemID == uuid.Nil {
			v.TimerItemID = uuid.New()
		}
		err := s.createTimerItem(v)
		if err != nil {
			s.logger.Error("error creating timerItem ", err)
			return err
		}
		break
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
		return blankItem, nil
	}
	messageItem, err := s.getMessageItemByID(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error querying messageItem ", err)
		return nil, err
	}
	if messageItem != nil {
		return messageItem, nil
	}
	speakerItem, err := s.getSpeakerItemByID(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error querying speakerItem ", err)
		return nil, err
	}
	if speakerItem != nil {
		return speakerItem, nil
	}
	lyricsItem, err := s.getLyricsItemByID(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error querying lyricsItem ", err)
		return nil, err
	}
	if lyricsItem != nil {
		return lyricsItem, nil
	}
	timerItem, err := s.getTimerItemByID(user.UserID, itemID)
	if err != nil {
		s.logger.Error("error querying timerItem ", err)
		return nil, err
	}
	if timerItem != nil {
		return timerItem, nil
	}
	return nil, nil
}

func (s lowerThirdsService) GetItems(ctx context.Context) (*[]entities.Item, error) {
	socialID := ctx.Value(helpers.SocialIDKey).(string)
	s.logger.Debug("GetItems for socialID ", socialID)
	user, err := s.GetUserBySocialID(ctx, socialID)
	if err != nil {
		s.logger.Error("User not found by socialID", err)
		return nil, err
	}
	s.logger.Debug("GetItems for userID ", user.UserID)

	// Query each type of item separately
	blankItems, err := s.getBlankItemsByUser(user.UserID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	messageItems, err := s.getMessageItemsByUser(user.UserID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	speakerItems, err := s.getSpeakerItemsByUser(user.UserID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	lyricsItems, err := s.getLyricsItemsByUser(user.UserID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	timerItems, err := s.getTimerItemsByUser(user.UserID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	// Combine in an Item interface set
	var allItems []entities.Item
	for _, b := range blankItems {
		allItems = append(allItems, b)
	}
	for _, m := range messageItems {
		allItems = append(allItems, m)
	}
	for _, s := range speakerItems {
		allItems = append(allItems, s)
	}
	for _, l := range lyricsItems {
		allItems = append(allItems, l)
	}
	for _, t := range timerItems {
		allItems = append(allItems, t)
	}

	// Sort items by GetOrder
	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].GetOrder() < allItems[j].GetOrder()
	})

	return &allItems, nil
}

func (s lowerThirdsService) GetItemsByMeeting(ctx context.Context, meetingID uuid.UUID) (*[]entities.Item, error) {
	socialID := ctx.Value(helpers.SocialIDKey).(string)
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
		allItems = append(allItems, b)
	}
	for _, m := range messageItems {
		allItems = append(allItems, m)
	}
	for _, s := range speakerItems {
		allItems = append(allItems, s)
	}
	for _, l := range lyricsItems {
		allItems = append(allItems, l)
	}
	for _, t := range timerItems {
		allItems = append(allItems, t)
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
		break
	case *entities.LyricsItem:
		if v.LyricsItemID == uuid.Nil {
			v.LyricsItemID = itemID
		}
		err := s.updateLyricsItem(itemID, v)
		if err != nil {
			s.logger.Error("error updating lyricsItem ", err)
			return err
		}
		break
	case *entities.MessageItem:
		if v.MessageItemID == uuid.Nil {
			v.MessageItemID = itemID
		}
		err := s.updateMessageItem(itemID, v)
		if err != nil {
			s.logger.Error("error updating messageItem ", err)
			return err
		}
		break
	case *entities.SpeakerItem:
		if v.SpeakerItemID == uuid.Nil {
			v.SpeakerItemID = itemID
		}
		err := s.updateSpeakerItem(itemID, v)
		if err != nil {
			s.logger.Error("error updating speakerItem ", err)
			return err
		}
		break
	case *entities.TimerItem:
		if v.TimerItemID == uuid.Nil {
			v.TimerItemID = itemID
		}
		err := s.updateTimerItem(itemID, v)
		if err != nil {
			s.logger.Error("error updating timerItem ", err)
			return err
		}
		break
	}
	return nil
}
