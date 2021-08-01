package middleware

import (
	"context"
	"net/http"
	"tacasa-web/options"

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

func setValueInContext(r *http.Request, key, value string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, key, value)
	return r.WithContext(ctx)
}
