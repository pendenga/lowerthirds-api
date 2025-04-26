package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"lowerthirdsapi/internal/helpers"
	"net/http"
	"strings"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

type Routes []Route

func (s *Server) Route() {
	var routes = Routes{
		// meetings
		Route{"getMeetings", "GET", "/meetings", s.getMeetings()},
		Route{"postMeeting", "POST", "/meetings", s.postMeeting()},
		Route{"getMeeting", "GET", "/meetings/{MeetingID}", s.getMeeting()},
		Route{"updateMeeting", "PUT", "/meetings/{MeetingID}", s.updateMeeting()},
		Route{"deleteMeeting", "DELETE", "/meetings/{MeetingID}", s.deleteMeeting()},
		Route{"getSlides", "GET", "/meetings/{MeetingID}/slides", s.getMeetingSlides()},

		// orgs
		Route{"getOrgs", "GET", "/orgs", s.getOrgs()},
		Route{"postOrg", "POST", "/orgs", s.postOrg()},
		Route{"getOrg", "GET", "/orgs/{OrgID}", s.getOrg()},
		Route{"updateOrg", "PUT", "/orgs/{OrgID}", s.updateOrg()},
		Route{"deleteOrg", "DELETE", "/orgs/{OrgID}", s.deleteOrg()},
		Route{"getOrgMeetings", "GET", "/orgs/{OrgID}/meetings", s.getOrgMeetings()},
		Route{"getOrgUsers", "GET", "/orgs/{OrgID}/users", s.getUsersByOrg()},

		// slides
		Route{"getSlides", "GET", "/slides", s.getSlides()},
		Route{"postSlide", "POST", "/slides", s.postSlide()},
		Route{"getSlide", "GET", "/slides/{SlideID}", s.getSlide()},
		Route{"updateSlide", "PUT", "/slides/{SlideID}", s.updateSlide()},
		Route{"deleteSlide", "DELETE", "/slides/{SlideID}", s.deleteSlide()},

		// users
		Route{"getUsers", "GET", "/users", s.getUsers()},
		Route{"postUser", "POST", "/users", s.postUser()},
		Route{"getUser", "GET", "/users/{UserID}", s.getUser()},
		Route{"updateUser", "PUT", "/users/{UserID}", s.updateUser()},
		Route{"deleteUser", "DELETE", "/users/{UserID}", s.deleteUser()},
		Route{"getUserMeetings", "GET", "/users/{UserID}/meetings", s.getUserMeetings()},
		Route{"getUserOrgs", "GET", "/users/{UserID}/orgs", s.getOrgsByUser()},
		Route{"setUserOrgs", "PUT", "/users/{UserID}/orgs", s.setOrgsByUser()},
	}
	for _, r := range routes {
		s.router.Use(authClaims(s.Logger))
		s.router.Handle(r.Pattern, r.Handler).Methods(r.Method).Name(r.Name)
	}

}

// authClaims is a middleware function to check auth headers
func authClaims(log *logrus.Entry) mux.MiddlewareFunc {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			w.Header().Set("Content-Type", "application/json")

			a := r.Header.Get("Authorization")
			if a == "" {
				// for testing
				ctx = context.WithValue(ctx, helpers.UserIDKey, uuid.MustParse("3cd5fe4e-9ecb-4ec2-b7c7-0d19288c08e0"))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			token := strings.TrimPrefix(a, "Bearer ")
			if token != "ABC123" {
				errMsg := "invalid token"
				log.Error(errMsg)
				http.Error(w, errMsg, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

// handleWithMiddleware wraps a route handler with any number of middleware functions
func handleWithMiddleware(handler http.Handler, mwf ...mux.MiddlewareFunc) http.Handler {
	for i := len(mwf) - 1; i >= 0; i-- {
		handler = mwf[i].Middleware(handler)
	}

	return handler
}
