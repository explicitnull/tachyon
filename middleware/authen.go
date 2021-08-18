package middleware

import (
	"html/template"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/sirupsen/logrus"
)

func (m *Middleware) CheckCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		le := logrus.
			WithField("origin", "middleware").
			WithField("requestID", ctx.Value("requestID"))

		username, ok := checkCookie(r, m.sc, le)
		if ok {
			le.WithField("username", username).Debugf("cookie verified")

			r = setValueInContext(r, "username", username)

			next.ServeHTTP(w, r)
		} else {
			le.Info("no cookie or bad cookie signature")

			t, err := template.ParseFiles("templates/login.htm")
			if err != nil {
				le.WithError(err).Error("template parsing failed")
				return
			}

			t.Execute(w, nil)
		}
	})
}

func checkCookie(r *http.Request, sc *securecookie.SecureCookie, le *logrus.Entry) (string, bool) {
	cookie, err := r.Cookie("username")
	if err != nil {
		return "", false
	}

	// TODO: log if cookie sig is invalid
	val := make(map[string]string)
	err = sc.Decode("username", cookie.Value, &val)
	if err != nil {
		le.WithError(err).Error("cookie decoding failed")
		return "", false
	}

	return val["name"], true
}
