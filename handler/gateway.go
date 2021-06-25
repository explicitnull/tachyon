package handler

import (
	"database/sql"
	"tac-gateway/options"
)

type Gateway struct {
	appName    string
	appVersion string
	Options    *options.Options
	db         *sql.DB
}

func NewGateway(o *options.Options, db *sql.DB) (*Gateway, error) {
	return &Gateway{
		appName:    "tac-gateway",
		appVersion: "1.0.0",
		Options:    o,
		db:         db,
	}, nil
}
