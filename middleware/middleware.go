package middleware

import (
	"tachyon-web/options"

	"github.com/gorilla/securecookie"
)

type Middleware struct {
	appName    string
	appVersion string
	Options    *options.Options
	sc         *securecookie.SecureCookie
}

func NewMiddleware(o *options.Options, sc *securecookie.SecureCookie) (*Middleware, error) {
	return &Middleware{
		Options: o,
		sc:      sc,
	}, nil
}
