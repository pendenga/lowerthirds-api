package server

import (
    "encoding/json"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "lowerthirdsapi/internal/entities"
    "lowerthirdsapi/internal/helpers"
    "net/http"
)

func (s *Server) deleteItem() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        itemID := uuid.MustParse(mux.Vars(req)["ItemID"])

        err := s.lowerThirdsService.DeleteItem(ctx, itemID)
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

        itemID := uuid.MustParse(mux.Vars(req)["ItemID"])

        item, err := s.lowerThirdsService.GetItem(ctx, itemID)
        if err != nil {
            s.Logger.Error(err)
            helpers.WriteError(ctx, err, w)
            return
        }

        if item == nil {
            s.Logger.Error("item not found ", itemID)
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
        ctx := req.Context()

        var item entities.Item

        if err := json.NewDecoder(req.Body).Decode(&item); err != nil {
            s.Logger.Error("Invalid JSON ", err)
            helpers.WriteError(ctx, err, w)
            return
        }

        err := s.lowerThirdsService.CreateItem(ctx, item)
        if err != nil {
            s.Logger.Error("CreateItem error ", err)
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

        itemID := uuid.MustParse(mux.Vars(req)["ItemID"])

        var item entities.Item
        var blankItem entities.BlankItem
        if err := json.NewDecoder(req.Body).Decode(&blankItem); err != nil {
            var lyricsItem entities.LyricsItem
            if err := json.NewDecoder(req.Body).Decode(&lyricsItem); err != nil {
                var messageItem entities.MessageItem
                if err := json.NewDecoder(req.Body).Decode(&messageItem); err != nil {
                    var speakerItem entities.SpeakerItem
                    if err := json.NewDecoder(req.Body).Decode(&speakerItem); err != nil {
                        var timerItem entities.TimerItem
                        if err := json.NewDecoder(req.Body).Decode(&timerItem); err != nil {
                            s.Logger.Error("Invalid JSON ", err)
                            helpers.WriteError(ctx, err, w)
                            return
                        } else {
                            item = timerItem
                        }
                    } else {
                        item = speakerItem
                    }
                } else {
                    item = messageItem
                }
            } else {
                item = lyricsItem
            }
        } else {
            item = blankItem
        }

        err := s.lowerThirdsService.UpdateItem(ctx, itemID, item)
        if err != nil {
            s.Logger.Error("UpdateItem error ", err)
            helpers.WriteError(ctx, err, w)
            return
        }

        w.WriteHeader(http.StatusOK)
        _ = json.NewEncoder(w).Encode(item)
    })
}
