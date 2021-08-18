package middleware

import (
	"net/http"
	"tacacs-webconsole/repository"

	"github.com/sirupsen/logrus"
)

func (m *Middleware) CheckAccessLevel(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		le := logrus.
			WithField("origin", "middleware").
			WithField("requestID", ctx.Value("requestID"))

		username, ok := ctx.Value("username").(string)
		if !ok {
			le.Warn("no username in context")
			return
		}

		level, err := repository.GetAccessLevel(le, m.asClient, username)
		if err != nil {
			le.WithError(err).Error("getting access level failure")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		switch level {
		case "level1":
			if (r.URL.Path == "/") || (r.URL.Path == "/myaccount/") || (r.URL.Path == "/lockout/") || (r.URL.Path == "/logout/") {
				next.ServeHTTP(w, r)
			} else {
				le.Warn("access forbidden")
				http.Error(w, "access forbidden", http.StatusForbidden)
			}
		case "level2":
			if r.URL.Path != "/settings/" {
				next.ServeHTTP(w, r)
			} else {
				le.Warn("access forbidden")
				http.Error(w, "access forbidden", http.StatusForbidden)
			}
		case "level3":
			next.ServeHTTP(w, r)
		default:
			le.Error("unknown access level received from database")
		}
	})
}
