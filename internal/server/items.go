package server

import (
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"lowerthirdsapi/internal/entities"
	"lowerthirdsapi/internal/helpers"
	"net/http"
)

func (s *Server) deleteItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		itemID, err := uuid.Parse(mux.Vars(req)["ItemID"])
		if err != nil {
			s.Logger.Error("[deleteItem] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		err = s.lowerThirdsService.DeleteItem(ctx, itemID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK) // 200 OK
	})
}

func (s *Server) getItems() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		items, err := s.lowerThirdsService.GetItems(ctx)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(items)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) getItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		itemID, err := uuid.Parse(mux.Vars(req)["ItemID"])
		if err != nil {
			s.Logger.Error("[getItem] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		item, err := s.lowerThirdsService.GetItem(ctx, itemID)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
			return
		}

		if item == nil {
			s.Logger.Error("[getItem] item not found ", itemID)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(item)
		if err != nil {
			s.Logger.Error(err)
			helpers.WriteError(ctx, err, w)
		}
	})
}

func (s *Server) postItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			s.Logger.Error("[postItem] Failed to read body: ", err)
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}
		s.Logger.Debug("[postItem] Body: ", string(bodyBytes))

		item, err := entities.ParseItemJSON(bodyBytes)
		if err != nil {
			if err == entities.ErrUnknownItemType {
				http.Error(w, "Unknown item type", http.StatusBadRequest)
			} else {
				helpers.WriteError(req.Context(), err, w)
			}
			return
		}

		ctx := req.Context()
		err = s.lowerThirdsService.CreateItem(ctx, item)
		if err != nil {
			// Check for MySQL duplicate entry error
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				http.Error(w, "[postItem] already exists", http.StatusConflict)
				return
			}
			s.Logger.Error("[postItem] CreateItem error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(item)
	})
}

func (s *Server) updateItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		itemID, err := uuid.Parse(mux.Vars(req)["ItemID"])
		if err != nil {
			s.Logger.Error("[updateItem] error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			s.Logger.Error("[updateItem] Failed to read body: ", err)
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}
		s.Logger.Debug("[updateItem] Body: ", string(bodyBytes))

		item, err := entities.ParseItemJSON(bodyBytes)
		if err != nil {
			if err == entities.ErrUnknownItemType {
				http.Error(w, "Unknown item type", http.StatusBadRequest)
			} else {
				helpers.WriteError(ctx, err, w)
			}
			return
		}

		err = s.lowerThirdsService.UpdateItem(ctx, itemID, item)
		if err != nil {
			// Check for MySQL duplicate entry error
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				http.Error(w, "[updateItem] already exists", http.StatusConflict)
				return
			}
			s.Logger.Error("[updateItem] UpdateItem error ", err)
			helpers.WriteError(ctx, err, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(item)
	})
}
