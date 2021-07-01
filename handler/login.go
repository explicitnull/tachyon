package handler

import (
	"bufio"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"strings"
	"tachyon-web/repository"
	"time"

	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
)

func (g *Gateway) Login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.htm")
	if err != nil {
		log.Errorf("template parsing failed: %v", err)
	}
	t.Execute(w, nil)
}

func (g *Gateway) LoginDo(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	ok := loginDo(username, password, g.db, r)
	if !ok {
		t, err := template.ParseFiles("templates/loginerror.htm")
		if err != nil {
			log.Errorf("%v", err)
			return
		}
		t.Execute(w, nil)

		return
	}

	err := setCookie(w, username, g.sc)
	if err != nil {
		fmt.Fprintf(w, "setting cookie failed")
		log.Errorf("%v", err)
		return
	}

	mid, err := template.ParseFiles("templates/loginok.htm")
	if err != nil {
		log.Errorf("%v", err)
		return
	}
	mid.Execute(w, nil)

}

func loginDo(username, formPassword string, db *sql.DB, r *http.Request) bool {
	ctx := r.Context()
	dbhash, err := repository.GetPasswordHash(db, username)
	if err != nil {
		log.Errorf("GetPasswordHash() failed: %v", err)
		return false
	}

	if dbhash == "" {
		log.WithField("username", username).WithField("ip", r.RemoteAddr).Warning("user not found")
		return false
	}

	hashParts := strings.Split(dbhash, "$")
	if len(hashParts) != 3 {
		log.WithField("username", username).Error("wrong database hash format")
		return false
	}

	salt := hashParts[2]
	formHash := hashPassword(salt, formPassword)

	if formHash != dbhash {
		log.WithField("username", username).WithField("ip", r.RemoteAddr).Warning("wrong password")
		return false
	}

	log.WithField("requestID", ctx.Value("requestID")).
		WithField("username", username).
		WithField("ip", r.RemoteAddr).
		Info("logged in")

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

func hashPassword(salt, password string) string {
	cmd := exec.Command("openssl", "passwd", "-1", "-salt", salt, password)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf("%v", err)
		return ""
	}

	cmd.Start()
	pipe := bufio.NewReader(stdout)

	line, _, err := pipe.ReadLine()
	if err != nil {
		log.Errorf("%v", err)
		return ""
	}

	return string(line)
}

