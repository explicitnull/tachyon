package middleware

import (
	"context"
	"net/http"
	"tachyon/options"

	"github.com/aerospike/aerospike-client-go"
	"github.com/gorilla/securecookie"
)

type Middleware struct {
	appName    string
	appVersion string
	Options    *options.Options
	sc         *securecookie.SecureCookie
	asClient   *aerospike.Client
}

func NewMiddleware(o *options.Options, sc *securecookie.SecureCookie, asClient *aerospike.Client) (*Middleware, error) {
	return &Middleware{
		Options:  o,
		sc:       sc,
		asClient: asClient,
	}, nil
}

func setValueInContext(r *http.Request, key, value string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, key, value)
	return r.WithContext(ctx)
}
