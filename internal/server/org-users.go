package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"lowerthirdsapi/internal/helpers"
	"net/http"
)

func (s *Server) getUsersByOrg() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		orgID := uuid.MustParse(mux.Vars(req)["OrgID"])

		meetings, err := s.lowerThirdsService.GetUsersByOrg(ctx, orgID)
		if err != nil {
			s.Logger.Error("GetUsersByOrg error ", err)
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

func (s *Server) getOrgsByUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		userID := mux.Vars(req)["UserID"]

		meetings, err := s.lowerThirdsService.GetOrgsByUser(ctx, userID)
		if err != nil {
			s.Logger.Error("GetOrgsByUser error ", err)
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

func (s *Server) setOrgsByUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		userID := mux.Vars(req)["UserID"]

		var orgIDs []uuid.UUID
		if err := json.NewDecoder(req.Body).Decode(&orgIDs); err != nil {
			s.Logger.Error("Invalid JSON ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		err := s.lowerThirdsService.SetOrgsByUser(ctx, userID, orgIDs)
		if err != nil {
			s.Logger.Error("SetOrgsByUser error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK) // 200 OK
	})
}
