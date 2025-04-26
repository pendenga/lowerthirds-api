package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
	"net/http"
)

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	FullName  string `json:"full_name"`
	LastName  string `json:"last_name"`
}

func (s *Server) deleteUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		userID := uuid.MustParse(mux.Vars(req)["UserID"])

		err := s.lowerThirdsService.DeleteUser(ctx, userID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK) // 200 OK
	})
}

func (s *Server) getUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		users, err := s.lowerThirdsService.GetUsers(ctx)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) getUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		userID := uuid.MustParse(mux.Vars(req)["UserID"])

		users, err := s.lowerThirdsService.GetUser(ctx, userID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) getUserMeetings() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		userID := uuid.MustParse(mux.Vars(req)["UserID"])

		meetings, err := s.lowerThirdsService.GetMeetingsByUser(ctx, userID)
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

func (s *Server) postUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var user entities.User

		if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
			s.Logger.Error("Invalid JSON ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		// If ID is not provided, generate a new one
		if user.UserID == uuid.Nil {
			user.UserID = uuid.New()
		}

		err := s.lowerThirdsService.CreateUser(ctx, &user)
		if err != nil {
			s.Logger.Error("CreateUser error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(user)
	})
}

func (s *Server) updateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		userID := uuid.MustParse(mux.Vars(req)["UserID"])
		var user entities.User

		if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
			s.Logger.Error("Invalid JSON ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		// Ensure the user ID from the path matches the payload and allow for exclusion in the payload
		user.UserID = userID

		err := s.lowerThirdsService.UpdateUser(ctx, userID, &user)
		if err != nil {
			s.Logger.Error("UpdateUser error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	})
}
