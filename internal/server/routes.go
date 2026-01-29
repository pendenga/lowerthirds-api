package server

import (
    "github.com/gorilla/mux"
    "net/http"
)

type Route struct {
    Name    string
    Method  string
    Pattern string
    Handler http.Handler
}
type Routes []Route

func (s *Server) Route() {
    s.Router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        s.Logger.Debug("Got a global OPTIONS request")
        w.WriteHeader(http.StatusOK)
    })

    // add middleware for every request
    s.Router.Use(authClaims(s.Logger))
    s.Router.Use(queryParametersInContext(s.Logger))

    var routes = Routes{
        // meetings
        Route{"getMeetings", "GET", "/v1/meetings", s.getMeetings()}, // need this? Should filter by org first
        Route{"postMeeting", "POST", "/v1/meetings", s.postMeeting()},
        Route{"getMeeting", "GET", "/v1/meetings/{MeetingID}", s.getMeeting()},
        Route{"updateMeeting", "PUT", "/v1/meetings/{MeetingID}", s.updateMeeting()},
        Route{"deleteMeeting", "DELETE", "/v1/meetings/{MeetingID}", s.deleteMeeting()},
        Route{"getMeetingItems", "GET", "/v1/meetings/{MeetingID}/items", s.getMeetingItems()}, // need this? Items are included in meeting

        // orgs
        Route{"getOrgs", "GET", "/v1/orgs", s.getOrgs()},
        Route{"postOrg", "POST", "/v1/orgs", s.postOrg()},
        Route{"getOrg", "GET", "/v1/orgs/{OrgID}", s.getOrg()},
        Route{"updateOrg", "PUT", "/v1/orgs/{OrgID}", s.updateOrg()},
        Route{"deleteOrg", "DELETE", "/v1/orgs/{OrgID}", s.deleteOrg()},
        Route{"getOrgMeetings", "GET", "/v1/orgs/{OrgID}/meetings", s.getOrgMeetings()},
        Route{"getOrgUsers", "GET", "/v1/orgs/{OrgID}/users", s.getUsersByOrg()},

        // items
        Route{"getItems", "GET", "/v1/items", s.getItems()},
        Route{"postItem", "POST", "/v1/items", s.postItem()},
        Route{"updateItem", "PUT", "/v1/items/{ItemID}", s.updateItem()},
        Route{"deleteItem", "DELETE", "/v1/items/{ItemID}", s.deleteItem()},

        // users
        Route{"getUsers", "GET", "/v1/users", s.getUsers()},
        Route{"postUser", "POST", "/v1/users", s.postUser()},
        Route{"getUser", "GET", "/v1/users/{UserID}", s.getUser()},
        Route{"updateUser", "PUT", "/v1/users/{UserID}", s.updateUser()},
        Route{"deleteUser", "DELETE", "/v1/users/{UserID}", s.deleteUser()},
        Route{"getUserMeetings", "GET", "/v1/users/{UserID}/meetings", s.getUserMeetings()}, // need this? Should go through org
        Route{"getUserOrgs", "GET", "/v1/users/{UserID}/orgs", s.getOrgsByUser()},
        Route{"setUserOrgs", "PUT", "/v1/users/{UserID}/orgs", s.setOrgsByUser()},
    }
    for _, r := range routes {
        s.Router.Handle(r.Pattern, r.Handler).Methods(r.Method).Name(r.Name)
    }
}

// handleWithMiddleware wraps a route handler with any number of middleware functions
func handleWithMiddleware(handler http.Handler, mwf ...mux.MiddlewareFunc) http.Handler {
    for i := len(mwf) - 1; i >= 0; i-- {
        handler = mwf[i].Middleware(handler)
    }

    return handler
}
