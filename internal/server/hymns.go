package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
	"net/http"
)

func (s *Server) deleteHymn() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		hymnID := uuid.MustParse(mux.Vars(req)["HymnID"])

		err := s.lowerThirdsService.DeleteHymn(ctx, hymnID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK) // 200 OK
	})
}

func (s *Server) getHymns() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		hymns, err := s.lowerThirdsService.GetHymns(ctx)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(hymns)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) getHymn() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		hymnID := uuid.MustParse(mux.Vars(req)["HymnID"])

		hymns, err := s.lowerThirdsService.GetHymn(ctx, hymnID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(hymns)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) postHymn() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var hymn entities.Hymn

		if err := json.NewDecoder(req.Body).Decode(&hymn); err != nil {
			s.Logger.Error("Invalid JSON ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		// If ID is not provided, generate a new one
		if hymn.HymnID == uuid.Nil {
			hymn.HymnID = uuid.New()
		}

		err := s.lowerThirdsService.CreateHymn(ctx, &hymn)
		if err != nil {
			s.Logger.Error("CreateHymn error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(hymn)
	})
}

func (s *Server) updateHymn() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		hymnID := uuid.MustParse(mux.Vars(req)["HymnID"])
		var hymn entities.Hymn

		if err := json.NewDecoder(req.Body).Decode(&hymn); err != nil {
			s.Logger.Error("Invalid JSON ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		// Ensure the hymn ID from the path matches the payload and allow for exclusion in the payload
		hymn.HymnID = hymnID

		err := s.lowerThirdsService.UpdateHymn(ctx, hymnID, &hymn)
		if err != nil {
			s.Logger.Error("UpdateHymn error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(hymn)
	})
}
