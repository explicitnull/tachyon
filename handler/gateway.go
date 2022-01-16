package handler

import (
	"tachyon/options"

	"github.com/aerospike/aerospike-client-go"

	"github.com/gorilla/securecookie"
)

type Gateway struct {
	appName         string
	appVersion      string
	Options         *options.Options
	aerospikeClient *aerospike.Client
	sc              *securecookie.SecureCookie
}

func NewGateway(o *options.Options, aerospikeClient *aerospike.Client, sc *securecookie.SecureCookie) (*Gateway, error) {
	return &Gateway{
		appName:         "tachyon",
		appVersion:      "1.0.0",
		Options:         o,
		aerospikeClient: aerospikeClient,
		sc:              sc,
	}, nil
}
