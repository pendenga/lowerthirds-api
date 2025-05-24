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
		Route{"getMeetings", "GET", "/meetings", s.getMeetings()}, // need this? Should filter by org first
		Route{"postMeeting", "POST", "/meetings", s.postMeeting()},
		Route{"getMeeting", "GET", "/meetings/{MeetingID}", s.getMeeting()},
		Route{"updateMeeting", "PUT", "/meetings/{MeetingID}", s.updateMeeting()},
		Route{"deleteMeeting", "DELETE", "/meetings/{MeetingID}", s.deleteMeeting()},
		Route{"getMeetingItems", "GET", "/meetings/{MeetingID}/items", s.getMeetingItems()}, // need this? Items are included in meeting

		// orgs
		Route{"getOrgs", "GET", "/orgs", s.getOrgs()},
		Route{"postOrg", "POST", "/orgs", s.postOrg()},
		Route{"getOrg", "GET", "/orgs/{OrgID}", s.getOrg()},
		Route{"updateOrg", "PUT", "/orgs/{OrgID}", s.updateOrg()},
		Route{"deleteOrg", "DELETE", "/orgs/{OrgID}", s.deleteOrg()},
		Route{"getOrgMeetings", "GET", "/orgs/{OrgID}/meetings", s.getOrgMeetings()},
		Route{"getOrgUsers", "GET", "/orgs/{OrgID}/users", s.getUsersByOrg()},

		// orgusers
		Route{"getOrgUserMap", "GET", "/orgusers", s.getOrgUsersMap()},
		Route{"postOrgUserMap", "POST", "/orgusers", s.postOrgUserMap()},
		Route{"deleteOrgUserMap", "DELETE", "/orgusers", s.deleteOrgUserMap()},

		// items
		Route{"getItems", "GET", "/items", s.getItems()}, // need this? Should just get items by meeting or by id
		Route{"postItem", "POST", "/items", s.postItem()},
		Route{"getItem", "GET", "/items/{ItemID}", s.getItem()},
		Route{"updateItem", "PUT", "/items/{ItemID}", s.updateItem()},
		Route{"deleteItem", "DELETE", "/items/{ItemID}", s.deleteItem()},

		// users
		Route{"getUsers", "GET", "/users", s.getUsers()},
		Route{"postUser", "POST", "/users", s.postUser()},
		Route{"getUser", "GET", "/users/{UserID}", s.getUser()},
		Route{"updateUser", "PUT", "/users/{UserID}", s.updateUser()},
		Route{"deleteUser", "DELETE", "/users/{UserID}", s.deleteUser()},
		Route{"getUserMeetings", "GET", "/users/{UserID}/meetings", s.getUserMeetings()}, // need this? Should go through org
		Route{"getUserOrgs", "GET", "/users/{UserID}/orgs", s.getOrgsByUser()},
		Route{"setUserOrgs", "PUT", "/users/{UserID}/orgs", s.setOrgsByUser()},
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
