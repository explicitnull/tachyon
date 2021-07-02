package handler

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const timeShort = "2006-01-02 15:04"

type Header struct {
	Name   string
	Item1  string
	Item2  string
	Item3  string
	Item4  string
	Item5  string
	Item6  string
	Item7  string
	Item8  string
	Item9  string
	Item10 string
}

func makeContextAndLogrusEntry(r *http.Request) (context.Context, *logrus.Entry) {
	ctx := r.Context()
	return ctx, log.WithField("requestID", ctx.Value("requestID")).WithField("username", ctx.Value("username"))
}
