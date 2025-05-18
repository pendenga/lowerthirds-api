package server

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
	"net/http"
)

func (s *Server) deleteMeeting() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		meetingID := uuid.MustParse(mux.Vars(req)["MeetingID"])

		err := s.lowerThirdsService.DeleteMeeting(ctx, meetingID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK) // 200 OK
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

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(meetings)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) getMeeting() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		meetingID := uuid.MustParse(mux.Vars(req)["MeetingID"])

		meetings, err := s.lowerThirdsService.GetMeeting(ctx, meetingID)
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

func (s *Server) getMeetingItems() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		meetingID := uuid.MustParse(mux.Vars(req)["MeetingID"])

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
			s.Logger.Error("Invalid JSON ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		if meeting.MeetingDate.IsZero() {
			s.Logger.Error("Missing or invalid meeting_date")
			helpers.WriteError(ctx, fmt.Errorf("date is required in RFC3339 format"), w)
			return
		}

		// If ID is not provided, generate a new one
		if meeting.MeetingID == uuid.Nil {
			meeting.MeetingID = uuid.New()
		}

		err := s.lowerThirdsService.CreateMeeting(ctx, &meeting)
		if err != nil {
			s.Logger.Error("CreateMeeting error ", err)
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

		meetingID := uuid.MustParse(mux.Vars(req)["MeetingID"])
		var meeting entities.Meeting

		if err := json.NewDecoder(req.Body).Decode(&meeting); err != nil {
			s.Logger.Error("Invalid JSON ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		if meeting.MeetingDate.IsZero() {
			s.Logger.Error("Missing or invalid meeting_date")
			helpers.WriteError(ctx, fmt.Errorf("date is required in RFC3339 format"), w)
			return
		}

		// Ensure the meeting ID from the path matches the payload and allow for exclusion in the payload
		meeting.MeetingID = meetingID

		err := s.lowerThirdsService.UpdateMeeting(ctx, meetingID, &meeting)
		if err != nil {
			s.Logger.Error("UpdateMeeting error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(meeting)
	})
}
