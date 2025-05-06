package server

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"lowerthirdsapi/internal/config"
	"lowerthirdsapi/internal/storage"
	"net/http"
	"time"
)

type Server struct {
	*http.Server
	DB                 *sqlx.DB
	Router             *mux.Router
	lowerThirdsService storage.LowerThirdsService
	Logger             *logrus.Entry
}

func New(cfg *config.Config, db *sqlx.DB, lowerThirdsService storage.LowerThirdsService, log *logrus.Entry) *Server {
	timeout := 20 * time.Second
	router := mux.NewRouter(mux.WithServiceName("lowerthirds-api"))
	router.StrictSlash(true)

	server := &Server{
		Server: &http.Server{
			Handler:        router,
			Addr:           ":9090",
			ReadTimeout:    timeout,
			WriteTimeout:   timeout,
			MaxHeaderBytes: 8192,
		},
		DB:                 db,
		lowerThirdsService: lowerThirdsService,
		Router:             router,
		Logger:             log,
	}
	server.Route()
	return server
}
