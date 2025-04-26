package main

import (
	"context"
	"github.com/sirupsen/logrus"
	ddsqlx "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
	"lowerthirdsapi/internal/config"
	"lowerthirdsapi/internal/helpers"
	"lowerthirdsapi/internal/logger"
	"lowerthirdsapi/internal/server"
	"lowerthirdsapi/internal/storage"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var log = logger.New()
	// Setup context that will cancel on signalled termination
	ctx := helpers.GetOsSignalContext(log)

	startApp(ctx, log)
}

func startApp(ctx context.Context, log *logrus.Entry) {
	log.Info("Starting up LowerThirds API")
	cfg := config.New(os.Getenv("ENV_FILES_DIR"))

	db := ddsqlx.MustConnect("mysql", cfg.MySQLConfig.ConnectionString())
	db.SetMaxOpenConns(cfg.MySQLConfig.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MySQLConfig.MaxOpenConns)
	db.SetConnMaxLifetime(time.Hour)
	defer db.Close()

	lowerThirdsService := storage.New(db, log)

	srvr := server.New(cfg, db, lowerThirdsService, log)
	defer helpers.ShutdownServer(srvr, log)
	go helpers.RunServer(srvr, log)

	// Exit safely
	<-ctx.Done()
	log.Info("exiting")
}
