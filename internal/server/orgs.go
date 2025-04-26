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

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(meetings)
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
			s.Logger.Error("Invalid JSON ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		// If ID is not provided, generate a new one
		if org.OrgID == uuid.Nil {
			org.OrgID = uuid.New()
		}

		err := s.lowerThirdsService.CreateOrg(ctx, &org)
		if err != nil {
			s.Logger.Error("CreateOrg error ", err)
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
			s.Logger.Error("Invalid JSON ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		// Ensure the org ID from the path matches the payload and allow for exclusion in the payload
		org.OrgID = orgID

		err := s.lowerThirdsService.UpdateOrg(ctx, orgID, &org)
		if err != nil {
			s.Logger.Error("UpdateOrg error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(org)
	})
}
