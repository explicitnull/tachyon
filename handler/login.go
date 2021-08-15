package handler

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"tacacs-webconsole/applogic"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/sirupsen/logrus"
)

const cookieTTL = 1440

func (g *Gateway) Login(w http.ResponseWriter, r *http.Request) {
	le := getLoggerWithoutUsername(r)

	t, err := template.ParseFiles("templates/login.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
	}
	t.Execute(w, nil)

	le.Info("handled ok")
}

func (g *Gateway) LoginAction(w http.ResponseWriter, r *http.Request) {
	le := getLoggerWithoutUsername(r)

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	le = le.WithField("username", username)

	ok := applogic.LoginAction(le, username, password, g.aerospikeClient)
	if !ok {
		t, err := template.ParseFiles("templates/loginerror.htm")
		if err != nil {
			le.WithError(err).Error("template parsing failed")
			http.Error(w, "template parsing failed", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)

		return
	}

	err := setCookie(w, username, g.sc)
	if err != nil {
		fmt.Fprintf(w, "setting cookie failed")
		le.WithError(err).Error("setting cookie failed")
		http.Error(w, "setting cookie failed", http.StatusFailedDependency)
		return
	}

	mid, err := template.ParseFiles("templates/loginok.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
		http.Error(w, "template parsing failed", http.StatusInternalServerError)
		return
	}
	mid.Execute(w, nil)

	le.Info("handled ok")
}

func setCookie(w http.ResponseWriter, username string, sc *securecookie.SecureCookie) error {
	value := map[string]string{
		"name": username,
	}

	expiration := time.Now().Add(cookieTTL * time.Minute)

	if encoded, err := sc.Encode("username", value); err == nil {
		cookie := &http.Cookie{
			Name:    "username",
			Value:   encoded,
			Path:    "/",
			Expires: expiration,
		}

		http.SetCookie(w, cookie)
	}
	return nil
}

func makeHashWithSalt(le *logrus.Entry, salt, password string) string {
	cmd := exec.Command("openssl", "passwd", "-1", "-salt", salt, password)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		le.Errorf("%v", err)
		return ""
	}

	cmd.Start()
	pipe := bufio.NewReader(stdout)

	line, _, err := pipe.ReadLine()
	if err != nil {
		le.Errorf("%v", err)
		return ""
	}

	return string(line)
}
