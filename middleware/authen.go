package middleware

import (
	"html/template"
	"net/http"

	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
)

func (m *Middleware) CheckCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		le := log.WithField("origin", "middleware").WithField("ip", r.RemoteAddr)
		user, ok := checkCookie(r, m.sc, le)
		if ok {
			le.WithField("username", user).Debugf("cookie verified")
			next.ServeHTTP(w, r)
		} else {
			t, err := template.ParseFiles("templates/login.htm")
			if err != nil {
				le.Errorf("template parsing failed: %v", err)
				return
			}
			t.Execute(w, nil)
		}
	})
}

func checkCookie(r *http.Request, sc *securecookie.SecureCookie, le log.Entry) (string, bool) {
	cookie, err := r.Cookie("username")
	switch err {
	case http.ErrNoCookie:
	case nil:
		value := make(map[string]string)
		if err = sc.Decode("username", cookie.Value, &value); err == nil {
			user := value["name"]
			return user, true
		}
	default:
		le.Error("cookie decoding failed")
	}

	return "", false
}
