package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
	"net/http"
)

func (s *Server) deleteMeeting() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		meetingID, err := uuid.Parse(mux.Vars(req)["MeetingID"])
		if err != nil {
			s.Logger.Error("[deleteMeeting] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		err = s.lowerThirdsService.DeleteMeeting(ctx, meetingID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204 No Content
	})
}

func (s *Server) getMeetings() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		meetings, err := s.lowerThirdsService.GetMeetings(ctx)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		// Get agenda items for every meeting
		// TODO: this method should be deprecated. The app should query for items by meeting
		var meetingsWithItems []entities.Meeting
		for _, meeting := range *meetings {
			agendaItems, err := s.lowerThirdsService.GetItemsByMeeting(ctx, meeting.MeetingID)
			if err != nil {
				s.Logger.Error(err)
				helpers.WriteError(ctx, err, w)
				return
			}

			meeting.AgendaItems = *agendaItems
			meetingsWithItems = append(meetingsWithItems, meeting)
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(meetingsWithItems)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) getMeeting() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		meetingID, err := uuid.Parse(mux.Vars(req)["MeetingID"])
		if err != nil {
			s.Logger.Error("[getMeeting] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		meeting, err := s.lowerThirdsService.GetMeeting(ctx, meetingID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		agendaItems, err := s.lowerThirdsService.GetItemsByMeeting(ctx, meetingID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		meeting.AgendaItems = *agendaItems

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(meeting)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) getMeetingItems() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		meetingID, err := uuid.Parse(mux.Vars(req)["MeetingID"])
		if err != nil {
			s.Logger.Error("[getMeetingItems] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		meetings, err := s.lowerThirdsService.GetItemsByMeeting(ctx, meetingID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(meetings)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) postMeeting() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var meeting entities.Meeting

		if err := json.NewDecoder(req.Body).Decode(&meeting); err != nil {
			s.Logger.Error("[postMeeting] ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		if meeting.MeetingDate.IsZero() {
			s.Logger.Error("[postMeeting] Missing or invalid meeting_date")
			helpers.WriteError(ctx, fmt.Errorf("date is required in RFC3339 format"), w)
			return
		}

		// If ID is not provided, generate a new one
		if meeting.MeetingID == uuid.Nil {
			meeting.MeetingID = uuid.New()
		}

		err := s.lowerThirdsService.CreateMeeting(ctx, &meeting)
		if err != nil {
			// Check for MySQL duplicate entry error
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				http.Error(w, "[postMeeting] already exists", http.StatusConflict)
				return
			}
			s.Logger.Error("[postMeeting] CreateMeeting error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(meeting)
	})
}

func (s *Server) updateMeeting() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		meetingID, err := uuid.Parse(mux.Vars(req)["MeetingID"])
		if err != nil {
			s.Logger.Error("[updateMeeting] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		var meeting entities.Meeting
		if err := json.NewDecoder(req.Body).Decode(&meeting); err != nil {
			s.Logger.Error("[updateMeeting] ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		if meeting.MeetingDate.IsZero() {
			s.Logger.Error("[updateMeeting] Missing or invalid meeting_date")
			helpers.WriteError(ctx, fmt.Errorf("date is required in RFC3339 format"), w)
			return
		}

		// Ensure the meeting ID from the path matches the payload and allow for exclusion in the payload
		meeting.MeetingID = meetingID

		err = s.lowerThirdsService.UpdateMeeting(ctx, meetingID, &meeting)
		if err != nil {
			// Check for MySQL duplicate entry error
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				http.Error(w, "[updateMeeting] already exists", http.StatusConflict)
				return
			}
			s.Logger.Error("[updateMeeting] UpdateMeeting error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(meeting)
	})
}
