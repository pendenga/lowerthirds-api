package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
	"net/http"
)

func (s *Server) deleteSlide() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		slideID := uuid.MustParse(mux.Vars(req)["SlideID"])

		err := s.lowerThirdsService.DeleteSlide(ctx, slideID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK) // 200 OK
	})
}

func (s *Server) getSlides() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		slides, err := s.lowerThirdsService.GetSlides(ctx)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(slides)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) getSlide() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		slideID := uuid.MustParse(mux.Vars(req)["SlideID"])

		slide, err := s.lowerThirdsService.GetSlide(ctx, slideID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		if slide == nil {
			s.Logger.Error("slide not found ", slideID)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(slide)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) postSlide() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var slide entities.Slide

		if err := json.NewDecoder(req.Body).Decode(&slide); err != nil {
			s.Logger.Error("Invalid JSON ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		err := s.lowerThirdsService.CreateSlide(ctx, slide)
		if err != nil {
			s.Logger.Error("CreateSlide error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(slide)
	})
}

func (s *Server) updateSlide() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		slideID := uuid.MustParse(mux.Vars(req)["SlideID"])

		var slide entities.Slide
		var blankSlide entities.BlankSlide
		if err := json.NewDecoder(req.Body).Decode(&blankSlide); err != nil {
			var lyricsSlide entities.LyricsSlide
			if err := json.NewDecoder(req.Body).Decode(&lyricsSlide); err != nil {
				var messageSlide entities.MessageSlide
				if err := json.NewDecoder(req.Body).Decode(&messageSlide); err != nil {
					var speakerSlide entities.SpeakerSlide
					if err := json.NewDecoder(req.Body).Decode(&speakerSlide); err != nil {
						var timerSlide entities.TimerSlide
						if err := json.NewDecoder(req.Body).Decode(&timerSlide); err != nil {
							s.Logger.Error("Invalid JSON ", err)
							helpers.WriteError(ctx, err, w)
							return
						} else {
							slide = timerSlide
						}
					} else {
						slide = speakerSlide
					}
				} else {
					slide = messageSlide
				}
			} else {
				slide = lyricsSlide
			}
		} else {
			slide = blankSlide
		}

		err := s.lowerThirdsService.UpdateSlide(ctx, slideID, slide)
		if err != nil {
			s.Logger.Error("UpdateSlide error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(slide)
	})
}
