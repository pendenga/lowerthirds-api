package helpers

import (
	"context"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func RunServer(srvr Server, log *logrus.Entry) {
	err := srvr.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func ShutdownServer(srvr Server, log *logrus.Entry) {
	err := srvr.Shutdown(context.Background())
	if err != nil {
		log.WithError(err).Error("failed to gracefully shutdown server")
	}
}
