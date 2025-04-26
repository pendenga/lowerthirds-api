package storage

import "fmt"

type MySQLConfig struct {
	Username     string `envconfig:"DB_USERNAME" default:"dba_lowerthirds"`
	Password     string `envconfig:"DB_PASSWORD" default:"3q.@JkTpbfFGifXGucqL"`
	Hostname     string `envconfig:"DB_HOSTNAME" default:"mysql.pendenga.com"`
	Schema       string `envconfig:"DB_SCHEMA" default:"lowerthirds"`
	Port         string `envconfig:"DB_PORT" default:"3306"`
	MaxOpenConns int    `envconfig:"DB_MAX_OPEN_CONNS" default:"10"`
}

func (cfg MySQLConfig) ConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		cfg.Username,
		cfg.Password,
		cfg.Hostname,
		cfg.Port,
		cfg.Schema,
	)
}
