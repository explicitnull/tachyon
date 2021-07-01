package handler

import (
	"bufio"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"regexp"

	log "github.com/sirupsen/logrus"
)

func (g *Gateway) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	le := log.
		WithField("requestID", ctx.Value("requestID")).
		WithField("username", ctx.Value("username"))

	header := Header{
		Name: ctx.Value("username").(string),
	}
	hdr, _ := template.ParseFiles("templates/hdr.htm")
	hdr.Execute(w, header)

	mid, _ := template.ParseFiles("templates/chpass.htm")
	mid.Execute(w, nil)

	ftr, _ := template.ParseFiles("templates/ftr.htm")
	ftr.Execute(w, nil)

	le.WithField("origin", "ChangePassword").Infof("request processed")
}

func (g *Gateway) ChangePasswordDo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	le := log.
		WithField("requestID", ctx.Value("requestID")).
		WithField("username", ctx.Value("username"))

	r.ParseForm()
	f1 := r.Form["pass1"]
	f2 := r.Form["pass2"]
	pass := f1[0]
	passConfirmation := f2[0]

	// checking if passwords don't match
	if pass != passConfirmation {
		fmt.Fprintf(w, "<p>Ошибка! Введенные пароли не совпадают.</p>")
		return
	}

	// checking if password has not [[:graph:]] symbols
	ok, _ := regexp.MatchString("^[[:graph:]]+$", pass)
	if !ok {
		fmt.Fprintln(w, "<p>Ошибка! Пароль содержит недопустимые символы. Забыли переключить раскладку?</p>")
		return
	}

	// checking password length
	CleanMap := make(map[string]interface{}, 0)
	CleanMap["pass"] = pass
	if len(pass) < g.Options.MinPassLen {
		fmt.Fprintln(w, "Ошибка! Пароль содержит слишком мало символов")
		return
	}

	// changing password
	username := ctx.Value("username").(string)

	err := updPass(g.db, username, pass)
	if err != nil {
		le.Errorf("%v", err)
		return
	}

	le.Info("user updated password")

	/*
		// Set flag "password changed"
		stmt, err2 := db.Prepare("update usr set pass_chd='true' where username=$1")
		checkErr(err2)
		defer stmt.Close()

		res, err3 := stmt.Exec(authUser)
		checkErr(err3)

		affect, err4 := res.RowsAffected()
		checkErr(err4)

		log.Println(affect, `rows changed while setting "password changed" flag for`, authUser)
	*/

	// activating user
	stmt, err := g.db.Prepare("update usr set act='true' where username=$1")
	if err != nil {
		le.Errorf("%v", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		le.Printf(`USER-4-ACTFAIL Error setting "active" flag for user %s\n`, username)
		return
	}

	le.Info("user status switched to active due to password change")

	fmt.Fprintf(w, "<p>Пароль изменен.</p>")
}

func updPass(db *sql.DB, usr, clText string) error {
	hash := makeHash(clText)

	stmt, err := db.Prepare("UPDATE usr SET pass=$1 WHERE username=$2")
	if err != nil {
		log.Errorf("%v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(hash, usr)
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	return err
}

func makeHash(pass string) string {
	/* This function generates and returns MD5 hashes for given passwords */
	cmd := exec.Command("openssl", "passwd", "-1", pass)
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

	hash := string(line)
	return hash
}
