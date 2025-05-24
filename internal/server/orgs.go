package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
	"net/http"
)

type OrgResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Server) deleteOrg() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		orgID := uuid.MustParse(mux.Vars(req)["OrgID"])

		err := s.lowerThirdsService.DeleteOrg(ctx, orgID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK) // 200 OK
	})
}

func (s *Server) getOrgs() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		orgs, err := s.lowerThirdsService.GetOrgs(ctx)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(orgs)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) getOrg() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		orgID := uuid.MustParse(mux.Vars(req)["OrgID"])

		orgs, err := s.lowerThirdsService.GetOrg(ctx, orgID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(orgs)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) getOrgMeetings() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		orgID := uuid.MustParse(mux.Vars(req)["OrgID"])

		meetings, err := s.lowerThirdsService.GetMeetingsByOrg(ctx, orgID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		// Get agenda items for each meeting
		// TODO: this is an inefficient way to get agenda items for each meeting
		// consider using a JOIN in the SQL query to get all items in one go
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

func (s *Server) postOrg() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var org entities.Organization

		if err := json.NewDecoder(req.Body).Decode(&org); err != nil {
			s.Logger.Error("[postOrg] ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		// If ID is not provided, generate a new one
		if org.OrgID == uuid.Nil {
			org.OrgID = uuid.New()
		}

		err := s.lowerThirdsService.CreateOrg(ctx, &org)
		if err != nil {
			s.Logger.Error("[postOrg] CreateOrg error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(org)
	})
}

func (s *Server) updateOrg() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		orgID := uuid.MustParse(mux.Vars(req)["OrgID"])
		var org entities.Organization

		if err := json.NewDecoder(req.Body).Decode(&org); err != nil {
			s.Logger.Error("[updateOrg] ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		// Ensure the org ID from the path matches the payload and allow for exclusion in the payload
		org.OrgID = orgID

		err := s.lowerThirdsService.UpdateOrg(ctx, orgID, &org)
		if err != nil {
			s.Logger.Error("[updateOrg] UpdateOrg error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(org)
	})
}
