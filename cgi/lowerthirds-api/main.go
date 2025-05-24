package main

import (
	ddsqlx "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
	"lowerthirdsapi/internal/config"
	"lowerthirdsapi/internal/logger"
	"lowerthirdsapi/internal/server"
	"lowerthirdsapi/internal/storage"
	"net/http/fcgi"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log := logger.New()
	cfg := config.New(os.Getenv("ENV_FILES_DIR"))

	db := ddsqlx.MustConnect("mysql", cfg.MySQLConfig.ConnectionString())
	db.SetMaxOpenConns(cfg.MySQLConfig.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MySQLConfig.MaxOpenConns)
	db.SetConnMaxLifetime(time.Hour)
	defer db.Close()

	lowerThirdsService := storage.New(db, log)
	srvr := server.New(cfg, db, lowerThirdsService, log)

	// Serve using FastCGI (or replace with cgi.Serve if you want classic CGI)
	if err := fcgi.Serve(nil, srvr.Router); err != nil {
		log.Fatal(err)
	}
}
