package middleware

import (
	"context"
	"html/template"
	"net/http"

	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
)

func (m *Middleware) CheckCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		le := log.
			WithField("origin", "middleware").
			WithField("requestID", ctx.Value("requestID"))

		username, ok := checkCookie(r, m.sc, le)
		if ok {
			le.WithField("username", username).Debugf("cookie verified")

			next.ServeHTTP(w, r)
		} else {
			t, err := template.ParseFiles("templates/login.htm")
			if err != nil {
				le.WithError(err).Error("template parsing failed")
				return
			}

			t.Execute(w, nil)
		}
	})
}

func checkCookie(r *http.Request, sc *securecookie.SecureCookie, le *log.Entry) (string, bool) {
	cookie, err := r.Cookie("username")
	if err != nil {
		return "", false
	}

	val := make(map[string]string)
	err = sc.Decode("username", cookie.Value, &val)
	if err != nil {
		le.Error("cookie decoding failed")
		return "", false
	}

	username := val["name"]
	setUsernameInContext(r, username)

	return username, true
}

func setUsernameInContext(r *http.Request, username string) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "username", username)
	r = r.WithContext(ctx)
}
