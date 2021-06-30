package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (m *Middleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		le := log.WithField("url", r.RequestURI).WithField("ip", r.RemoteAddr).WithField("origin", "middleware")

		ctx := r.Context()

		u := ctx.Value("username")
		if u == nil {
			le.Info("request received")
		} else {
			username, ok := u.(string)
			if !ok {
				log.Warnf("context username type error")
			}

			le.WithField("username", username).Info("request received")
		}
		next.ServeHTTP(w, r)
	})
}
