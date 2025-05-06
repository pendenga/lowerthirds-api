package storage

import (
	"context"
	"github.com/google/uuid"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
	"sort"
)

func (s lowerThirdsService) CreateSlide(ctx context.Context, slide entities.Slide) error {
	userID := ctx.Value(helpers.UserIDKey).(string)
	s.logger.Debug("CreateSlides for userID ", userID)

	// Each type of slide handled separately
	switch v := slide.(type) {
	case *entities.BlankSlide:
		if v.BlankSlideID == uuid.Nil {
			v.BlankSlideID = uuid.New()
		}
		err := s.createBlankSlide(v)
		if err != nil {
			s.logger.Error("error creating blankSlide ", err)
			return err
		}
		break
	case *entities.LyricsSlide:
		if v.LyricsSlideID == uuid.Nil {
			v.LyricsSlideID = uuid.New()
		}
		err := s.createLyricsSlide(v)
		if err != nil {
			s.logger.Error("error creating lyricsSlide ", err)
			return err
		}
		break
	case *entities.MessageSlide:
		if v.MessageSlideID == uuid.Nil {
			v.MessageSlideID = uuid.New()
		}
		err := s.createMessageSlide(v)
		if err != nil {
			s.logger.Error("error creating messageSlide ", err)
			return err
		}
		break
	case *entities.SpeakerSlide:
		if v.SpeakerSlideID == uuid.Nil {
			v.SpeakerSlideID = uuid.New()
		}
		err := s.createSpeakerSlide(v)
		if err != nil {
			s.logger.Error("error creating speakerSlide ", err)
			return err
		}
		break
	case *entities.TimerSlide:
		if v.TimerSlideID == uuid.Nil {
			v.TimerSlideID = uuid.New()
		}
		err := s.createTimerSlide(v)
		if err != nil {
			s.logger.Error("error creating timerSlide ", err)
			return err
		}
		break
	}
	return nil
}

func (s lowerThirdsService) DeleteSlide(ctx context.Context, slideID uuid.UUID) error {
	userID := ctx.Value(helpers.UserIDKey).(string)
	s.logger.Debug("DeleteSlides for userID ", userID)

	var totalAffectedRows int64 = 0

	// Query each type of slide separately
	affectedRows, err := s.deleteBlankSlide(userID, slideID)
	if err != nil {
		s.logger.Error("error deleting blankSlide ", err)
		return err
	}
	totalAffectedRows = totalAffectedRows + affectedRows
	affectedRows, err = s.deleteMessageSlide(userID, slideID)
	if err != nil {
		s.logger.Error("error deleting messageSlide ", err)
		return err
	}
	totalAffectedRows = totalAffectedRows + affectedRows
	affectedRows, err = s.deleteSpeakerSlide(userID, slideID)
	if err != nil {
		s.logger.Error("error deleting speakerSlide ", err)
		return err
	}
	totalAffectedRows = totalAffectedRows + affectedRows
	affectedRows, err = s.deleteLyricsSlide(userID, slideID)
	if err != nil {
		s.logger.Error("error deleting lyricsSlide ", err)
		return err
	}
	totalAffectedRows = totalAffectedRows + affectedRows
	affectedRows, err = s.deleteTimerSlide(userID, slideID)
	if err != nil {
		s.logger.Error("error deleting timerSlide ", err)
		return err
	}
	totalAffectedRows = totalAffectedRows + affectedRows

	s.logger.Info("DeleteSlides affectedRows rows: ", totalAffectedRows)

	return nil
}

func (s lowerThirdsService) GetSlide(ctx context.Context, slideID uuid.UUID) (entities.Slide, error) {
	userID := ctx.Value(helpers.UserIDKey).(string)
	s.logger.Debug("GetSlides for userID ", userID)

	// Query each type of slide separately
	blankSlide, err := s.getBlankSlideByID(userID, slideID)
	if err != nil {
		s.logger.Error("error querying blankSlide ", err)
		return nil, err
	}
	if blankSlide != nil {
		return blankSlide, nil
	}
	messageSlide, err := s.getMessageSlideByID(userID, slideID)
	if err != nil {
		s.logger.Error("error querying messageSlide ", err)
		return nil, err
	}
	if messageSlide != nil {
		return messageSlide, nil
	}
	speakerSlide, err := s.getSpeakerSlideByID(userID, slideID)
	if err != nil {
		s.logger.Error("error querying speakerSlide ", err)
		return nil, err
	}
	if speakerSlide != nil {
		return speakerSlide, nil
	}
	lyricsSlide, err := s.getLyricsSlideByID(userID, slideID)
	if err != nil {
		s.logger.Error("error querying lyricsSlide ", err)
		return nil, err
	}
	if lyricsSlide != nil {
		return lyricsSlide, nil
	}
	timerSlide, err := s.getTimerSlideByID(userID, slideID)
	if err != nil {
		s.logger.Error("error querying timerSlide ", err)
		return nil, err
	}
	if timerSlide != nil {
		return timerSlide, nil
	}
	return nil, nil
}

func (s lowerThirdsService) GetSlides(ctx context.Context) (*[]entities.Slide, error) {
	userID := ctx.Value(helpers.UserIDKey).(string)
	s.logger.Debug("GetSlides for userID ", userID)

	// Query each type of slide separately
	blankSlides, err := s.getBlankSlidesByUser(userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	messageSlides, err := s.getMessageSlidesByUser(userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	speakerSlides, err := s.getSpeakerSlidesByUser(userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	lyricsSlides, err := s.getLyricsSlidesByUser(userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	timerSlides, err := s.getTimerSlidesByUser(userID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	// Combine in a Slide interface set
	var allSlides []entities.Slide
	for _, b := range blankSlides {
		allSlides = append(allSlides, b)
	}
	for _, m := range messageSlides {
		allSlides = append(allSlides, m)
	}
	for _, s := range speakerSlides {
		allSlides = append(allSlides, s)
	}
	for _, l := range lyricsSlides {
		allSlides = append(allSlides, l)
	}
	for _, t := range timerSlides {
		allSlides = append(allSlides, t)
	}

	// Sort slides by GetOrder
	sort.Slice(allSlides, func(i, j int) bool {
		return allSlides[i].GetOrder() < allSlides[j].GetOrder()
	})

	return &allSlides, nil
}

func (s lowerThirdsService) GetSlidesByMeeting(ctx context.Context, meetingID uuid.UUID) (*[]entities.Slide, error) {
	userID := ctx.Value(helpers.UserIDKey).(string)
	s.logger.Debug("GetSlidesByMeeting for userID ", userID, ", meetingID ", meetingID)

	// Query each type of slide separately
	blankSlides, err := s.getBlankSlidesByMeeting(userID, meetingID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	messageSlides, err := s.getMessageSlidesByMeeting(userID, meetingID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	speakerSlides, err := s.getSpeakerSlidesByMeeting(userID, meetingID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	lyricsSlides, err := s.getLyricsSlidesByMeeting(userID, meetingID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	timerSlides, err := s.getTimerSlidesByMeeting(userID, meetingID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	// Combine in a Slide interface set
	var allSlides []entities.Slide
	for _, b := range blankSlides {
		allSlides = append(allSlides, b)
	}
	for _, m := range messageSlides {
		allSlides = append(allSlides, m)
	}
	for _, s := range speakerSlides {
		allSlides = append(allSlides, s)
	}
	for _, l := range lyricsSlides {
		allSlides = append(allSlides, l)
	}
	for _, t := range timerSlides {
		allSlides = append(allSlides, t)
	}

	// Sort slides by GetOrder
	sort.Slice(allSlides, func(i, j int) bool {
		return allSlides[i].GetOrder() < allSlides[j].GetOrder()
	})

	return &allSlides, nil
}

func (s lowerThirdsService) UpdateSlide(ctx context.Context, slideID uuid.UUID, slide entities.Slide) error {
	userID := ctx.Value(helpers.UserIDKey).(string)
	s.logger.Debug("UpdateSlides for userID ", userID)

	// TODO: somehow check that the ID being used matches up with an existing record of the same type.
	// Theoretically we could get duplicate ID numbers in different tables

	// Each type of slide handled separately
	switch v := slide.(type) {
	case *entities.BlankSlide:
		if v.BlankSlideID == uuid.Nil {
			v.BlankSlideID = slideID
		}
		err := s.updateBlankSlide(slideID, v)
		if err != nil {
			s.logger.Error("error updating blankSlide ", err)
			return err
		}
		break
	case *entities.LyricsSlide:
		if v.LyricsSlideID == uuid.Nil {
			v.LyricsSlideID = slideID
		}
		err := s.updateLyricsSlide(slideID, v)
		if err != nil {
			s.logger.Error("error updating lyricsSlide ", err)
			return err
		}
		break
	case *entities.MessageSlide:
		if v.MessageSlideID == uuid.Nil {
			v.MessageSlideID = slideID
		}
		err := s.updateMessageSlide(slideID, v)
		if err != nil {
			s.logger.Error("error updating messageSlide ", err)
			return err
		}
		break
	case *entities.SpeakerSlide:
		if v.SpeakerSlideID == uuid.Nil {
			v.SpeakerSlideID = slideID
		}
		err := s.updateSpeakerSlide(slideID, v)
		if err != nil {
			s.logger.Error("error updating speakerSlide ", err)
			return err
		}
		break
	case *entities.TimerSlide:
		if v.TimerSlideID == uuid.Nil {
			v.TimerSlideID = slideID
		}
		err := s.updateTimerSlide(slideID, v)
		if err != nil {
			s.logger.Error("error updating timerSlide ", err)
			return err
		}
		break
	}
	return nil
}
