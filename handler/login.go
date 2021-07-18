package handler

import (
	"bufio"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"tachyon-web/repository"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/sirupsen/logrus"
)

func (g *Gateway) Login(w http.ResponseWriter, r *http.Request) {
	le := getLoggerWithoutUsername(r)

	t, err := template.ParseFiles("templates/login.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
	}
	t.Execute(w, nil)
}

func (g *Gateway) LoginDo(w http.ResponseWriter, r *http.Request) {
	le := getLoggerWithoutUsername(r)

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	le = le.WithField("username", username)

	ok := loginDo(le, username, password, g.db, r)
	if !ok {
		t, err := template.ParseFiles("templates/loginerror.htm")
		if err != nil {
			le.Errorf("%v", err)
			return
		}
		t.Execute(w, nil)

		return
	}

	err := setCookie(w, username, g.sc)
	if err != nil {
		fmt.Fprintf(w, "setting cookie failed")
		le.Errorf("%v", err)
		return
	}

	mid, err := template.ParseFiles("templates/loginok.htm")
	if err != nil {
		le.Errorf("%v", err)
		return
	}
	mid.Execute(w, nil)

}

func loginDo(le *logrus.Entry, username, formPassword string, db *sql.DB, r *http.Request) bool {
	dbhash, err := repository.GetPasswordHash(db, username)
	if err != nil {
		le.WithError(err).Errorf("GetPasswordHash() failed")
		return false
	}

	if dbhash == "" {
		le.Warning("user not found")
		return false
	}

	// hashParts := strings.Split(dbhash, "$")
	// if len(hashParts) != 3 {
	// 	le.Error("wrong database hash format")
	// 	return false
	// }

	// salt := hashParts[2]
	formHash := makeHash(le, formPassword)

	if formHash != dbhash {
		le.Warning("wrong password")
		return false
	}

	le.Info("logged in")

	return true
}

func setCookie(w http.ResponseWriter, username string, sc *securecookie.SecureCookie) error {
	value := map[string]string{
		"name": username,
	}
	expiration := time.Now().Add(1 * time.Hour)
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

func hashPassword(le *logrus.Entry, salt, password string) string {
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
