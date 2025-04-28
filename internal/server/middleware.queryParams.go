package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"lowerthirdsapi/internal/helpers"
	"net/http"
	"strconv"
	"time"
)

// queryParametersInContext is a middleware function to parse query parameters and put them in the context
func queryParametersInContext(log *logrus.Entry) mux.MiddlewareFunc {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("queryParametersInContext middleware")

			qp := helpers.DefaultQueryParams()

			query := r.URL.Query()
			if pageStr := query.Get("Page"); pageStr != "" {
				if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
					qp.Page = p
				}
			}
			if pageSizeStr := query.Get("PageSize"); pageSizeStr != "" {
				if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
					qp.PageSize = ps
				}
			}
			if dateFromStr := query.Get("DateFrom"); dateFromStr != "" {
				df, err := time.Parse("2006-01-02", dateFromStr)
				if err != nil {
					errMsg := "invalid date format (DateFrom)"
					log.Error(errMsg)
					http.Error(w, errMsg, http.StatusBadRequest)
					return
				}
				qp.DateFrom = &df
			}
			if dateToStr := query.Get("DateTo"); dateToStr != "" {
				dt, err := time.Parse("2006-01-02", dateToStr)
				if err != nil {
					errMsg := "invalid date format (DateFrom)"
					log.Error(errMsg)
					http.Error(w, errMsg, http.StatusBadRequest)
					return
				}
				qp.DateTo = &dt
			}
			if languageStr := query.Get("Language"); languageStr != "" {
				qp.Language = languageStr
			}

			newContext := context.WithValue(r.Context(), helpers.QueryParametersKey, qp)
			next.ServeHTTP(w, r.WithContext(newContext))
		})
	})
}
