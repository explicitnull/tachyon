package handler

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"regexp"
	"tachyon-web/options"
	"tachyon-web/repository"

	"github.com/sirupsen/logrus"
)

var (
	passwordsMismatchError = errors.New("passwords mismatch")
	badCharactersError     = errors.New("bad characters in password")
	tooShortError          = errors.New("too short password")
)

func (g *Gateway) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

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
	le := getLogger(r)

	r.ParseForm()
	f1 := r.Form["pass1"]
	f2 := r.Form["pass2"]
	pass := f1[0]
	passConfirm := f2[0]

	err := changePasswordDo(le, g.db, pass, passConfirm, g.Options, ctx)
	if err == passwordsMismatchError {
		fmt.Fprintf(w, "<p>Ошибка! Введенные пароли не совпадают.</p>")
		return
	} else if err == badCharactersError {
		fmt.Fprintln(w, "<p>Ошибка! Пароль содержит недопустимые символы. Забыли переключить раскладку?</p>")
		return
	} else if err == tooShortError {
		fmt.Fprintln(w, "Ошибка! Пароль содержит слишком мало символов")
		return
	}

	fmt.Fprintf(w, "<p>Пароль изменен.</p>")
}

func changePasswordDo(le *logrus.Entry, db *sql.DB, pass, passConfirm string, o *options.Options, ctx context.Context) error {
	// checking if passwords don't match
	if pass != passConfirm {
		return passwordsMismatchError
	}

	// checking if password has not [[:graph:]] symbols
	ok, _ := regexp.MatchString("^[[:graph:]]+$", pass)
	if !ok {
		return badCharactersError
	}

	// checking password length
	CleanMap := make(map[string]interface{}, 0)
	CleanMap["pass"] = pass
	if len(pass) < o.MinPassLen {
		return tooShortError
	}

	// changing password
	username := ctx.Value("username").(string)

	hash := makeHash(le, pass)

	err := repository.UpdatePassword(db, username, hash)
	if err != nil {
		le.WithError(err).Error("password update failed")
		return err
	}

	le.Info("user password updated")

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
	err = repository.ActivateUser(db, le, username)
	if err != nil {
		le.WithError(err).Error("user activation failed")
		return err
	}

	le.Info("user status switched to active due to password update")

	return nil
}

/* makeHash generates MD5 hashes for given passwords */
func makeHash(le *logrus.Entry, pass string) string {
	cmd := exec.Command("openssl", "passwd", "-1", pass)
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

	hash := string(line)
	return hash
}
