package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
	"net/http"
	"strconv"
)

func (s *Server) deleteVerse() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		hymnID := uuid.MustParse(mux.Vars(req)["HymnID"])
		verseStr := mux.Vars(req)["VerseNum"]
		verseNum, err := strconv.Atoi(verseStr)
		if err != nil {
			s.Logger.Error("VerseNum needs to be integer ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		err = s.lowerThirdsService.DeleteVerse(ctx, hymnID, verseNum)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK) // 200 OK
	})
}

func (s *Server) getVerses() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		hymnID := uuid.MustParse(mux.Vars(req)["HymnID"])

		verses, err := s.lowerThirdsService.GetVerses(ctx, hymnID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(verses)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) getVerse() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		hymnID := uuid.MustParse(mux.Vars(req)["HymnID"])
		verseStr := mux.Vars(req)["VerseNum"]
		verseNum, err := strconv.Atoi(verseStr)
		if err != nil {
			s.Logger.Error("VerseNum needs to be integer ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		verses, err := s.lowerThirdsService.GetVerse(ctx, hymnID, verseNum)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(verses)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) postVerse() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		hymnID := uuid.MustParse(mux.Vars(req)["HymnID"])
		var verse entities.HymnVerse

		if err := json.NewDecoder(req.Body).Decode(&verse); err != nil {
			s.Logger.Error("Invalid JSON ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		err := s.lowerThirdsService.CreateVerse(ctx, hymnID, &verse)
		if err != nil {
			s.Logger.Error("CreateVerse error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(verse)
	})
}

func (s *Server) updateVerse() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		hymnID := uuid.MustParse(mux.Vars(req)["HymnID"])
		verseStr := mux.Vars(req)["VerseNum"]
		verseNum, err := strconv.Atoi(verseStr)
		if err != nil {
			s.Logger.Error("VerseNum needs to be integer ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		var verse entities.HymnVerse

		if err := json.NewDecoder(req.Body).Decode(&verse); err != nil {
			s.Logger.Error("Invalid JSON ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		err = s.lowerThirdsService.UpdateVerse(ctx, hymnID, verseNum, &verse)
		if err != nil {
			s.Logger.Error("UpdateVerse error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(verse)
	})
}
