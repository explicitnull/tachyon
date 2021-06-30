package handler

import (
	"database/sql"
	"tachyon-web/options"

	"github.com/gorilla/securecookie"
)

type Gateway struct {
	appName    string
	appVersion string
	Options    *options.Options
	db         *sql.DB
	sc         *securecookie.SecureCookie
}

func NewGateway(o *options.Options, db *sql.DB, sc *securecookie.SecureCookie) (*Gateway, error) {
	return &Gateway{
		appName:    "tachyon-web",
		appVersion: "1.0.0",
		Options:    o,
		db:         db,
		sc:         sc,
	}, nil
}
