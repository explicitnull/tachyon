package middleware

import (
	"context"
	"net/http"

	"github.com/dchest/uniuri"
	log "github.com/sirupsen/logrus"
)

func (m *Middleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := generateID()

		log.WithField("origin", "middleware").
			WithField("url", r.RequestURI).
			WithField("ip", r.RemoteAddr).
			WithField("requestID", requestID).
			Info("request received")

		setRequestIDInContext(r, requestID)

		next.ServeHTTP(w, r)
	})
}

func generateID() string {
	return uniuri.New()
}

func setRequestIDInContext(r *http.Request, requestID string) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "requestID", requestID)
	r = r.WithContext(ctx)
}
