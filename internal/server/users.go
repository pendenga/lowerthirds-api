package server

import (
	"encoding/json"
	"github.com/go-sql-driver/mysql"
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

		userID, err := uuid.Parse(mux.Vars(req)["UserID"])
		if err != nil {
			s.Logger.Error("[deleteUser] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		err = s.lowerThirdsService.DeleteUser(ctx, userID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204 No Content
	})
}

func (s *Server) getOrgsByUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		userID, err := uuid.Parse(mux.Vars(req)["UserID"])
		if err != nil {
			s.Logger.Error("[getOrgsByUser] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		meetings, err := s.lowerThirdsService.GetOrgsByUser(ctx, userID)
		if err != nil {
			s.Logger.Error("[getOrgsByUser] error ", err)
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

		userID, err := uuid.Parse(mux.Vars(req)["UserID"])
		if err != nil {
			s.Logger.Error("[getUser] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

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

		userID, err := uuid.Parse(mux.Vars(req)["UserID"])
		if err != nil {
			s.Logger.Error("[getUserMeetings] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

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
			s.Logger.Error("[postUser] ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		// If ID is not provided, generate a new one
		if user.UserID == uuid.Nil {
			user.UserID = uuid.New()
		}

		err := s.lowerThirdsService.CreateUser(ctx, &user)
		if err != nil {
			// Check for MySQL duplicate entry error
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				http.Error(w, "[postUser] already exists", http.StatusConflict)
				return
			}
			s.Logger.Error("[postUser] CreateUser error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(user)
	})
}

func (s *Server) setOrgsByUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		userID, err := uuid.Parse(mux.Vars(req)["UserID"])
		if err != nil {
			s.Logger.Error("[setOrgsByUser] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		var orgIDs []uuid.UUID
		if err := json.NewDecoder(req.Body).Decode(&orgIDs); err != nil {
			s.Logger.Error("[setOrgsByUser] ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		err = s.lowerThirdsService.SetOrgsByUser(ctx, userID, orgIDs)
		if err != nil {
			s.Logger.Error("[setOrgsByUser] SetOrgsByUser error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK) // 200 OK
	})
}

func (s *Server) updateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		userID, err := uuid.Parse(mux.Vars(req)["UserID"])
		if err != nil {
			s.Logger.Error("[updateUser] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		var user entities.User
		if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
			s.Logger.Error("[updateUser] ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		// Ensure the user ID from the path matches the payload and allow for exclusion in the payload
		user.UserID = userID

		err = s.lowerThirdsService.UpdateUser(ctx, userID, &user)
		if err != nil {
			// Check for MySQL duplicate entry error
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				http.Error(w, "[updateUser] already exists", http.StatusConflict)
				return
			}
			s.Logger.Error("[updateUser] UpdateUser error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	})
}
