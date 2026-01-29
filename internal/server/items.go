package server

import (
    "encoding/json"
    "io"
    "lowerthirdsapi/internal/entities"
    "lowerthirdsapi/internal/helpers"
    "net/http"

    "github.com/go-sql-driver/mysql"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
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

        w.WriteHeader(http.StatusNoContent) // 204 No Content
    })
}

func (s *Server) getItems() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        items, err := s.lowerThirdsService.GetItems(ctx)
        if err != nil {
            s.Logger.Error("[getItems] error ", err)
            helpers.WriteError(ctx, err, w)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        _ = json.NewEncoder(w).Encode(items)
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
