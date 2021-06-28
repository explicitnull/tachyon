package middleware

import (
	"html/template"
	"net/http"

	"github.com/gorilla/securecookie"

	log "github.com/sirupsen/logrus"
)

var hashKey = []byte("secret")
var blockKey = []byte("1234567890123456")
var sc = securecookie.New(hashKey, blockKey)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u := ctx.Value("username")
		if u == nil {
			log.Errorf("context username not set")
		}

		username, ok := u.(string)
		if !ok {
			log.Errorf("context username type error")
		}

		log.WithField("url", r.RequestURI).WithField("username", username).Info("request received")
		next.ServeHTTP(w, r)
	})
}

func CheckCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := checkCookie(r)
		if ok {
			log.WithField("username", user).Debugf("cookie verified")
			next.ServeHTTP(w, r)
		} else {
			t, err := template.ParseFiles("templates/login.htm")
			if err != nil {
				log.Errorf("template parsing failed: %v", err)
				return
			}
			t.Execute(w, nil)
		}
	})
}

func checkCookie(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("username")
	if err == nil {
		value := make(map[string]string)
		if err = sc.Decode("username", cookie.Value, &value); err == nil {
			user := value["name"]
			return user, true
		}
	}

	return "", false
}
