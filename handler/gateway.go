package handler

import (
	"database/sql"
	"tachyon-web/options"
)

type Gateway struct {
	appName    string
	appVersion string
	Options    *options.Options
	db         *sql.DB
}

func NewGateway(o *options.Options, db *sql.DB) (*Gateway, error) {
	return &Gateway{
		appName:    "tachyon-web",
		appVersion: "1.0.0",
		Options:    o,
		db:         db,
	}, nil
}
