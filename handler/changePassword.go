package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"tachyon-web/options"
	"tachyon-web/repository"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

var (
	errPasswordsMismatch = errors.New("passwords mismatch")
	badCharactersError   = errors.New("bad characters in password")
	tooShortError        = errors.New("too short password")
)

func (g *Gateway) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		fmt.Fprintf(w, "access forbidden")
		return
	}

	executeHeaderTemplate(le, w, username)

	mid, _ := template.ParseFiles("templates/chpass.htm")
	mid.Execute(w, nil)

	executeFooterTemplate(le, w)

	le.WithField("origin", "ChangePassword").Infof("request processed")
}

func (g *Gateway) ChangePasswordDo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		fmt.Fprintf(w, "access forbidden")
		return
	}

	r.ParseForm()
	f1 := r.Form["pass1"]
	f2 := r.Form["pass2"]
	pass := f1[0]
	passConfirm := f2[0]

	err := changePasswordDo(le, g.aerospikeClient, authenticatedUsername, pass, passConfirm, g.Options)
	if err == errPasswordsMismatch {
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

func changePasswordDo(le *logrus.Entry, aClient *aerospike.Client, authenticatedUsername, pass, passConfirm string, o *options.Options) error {
	// checking if passwords don't match
	if pass != passConfirm {
		return errPasswordsMismatch
	}

	// checking if password has not [[:graph:]] symbols
	ok, _ := regexp.MatchString("^[[:graph:]]+$", pass)
	if !ok {
		return badCharactersError
	}

	// checking password length
	CleanMap := make(map[string]interface{})
	CleanMap["pass"] = pass
	if len(pass) < o.MinPassLen {
		return tooShortError
	}

	// changing password
	hash := makeHash(le, pass)

	err := repository.UpdatePassword(aClient, authenticatedUsername, hash)
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
	err = repository.SetUserStatus(le, authenticatedUsername, true)
	if err != nil {
		le.WithError(err).Error("user activation failed")
		return err
	}

	le.Info("user status switched to active due to password update")

	return nil
}

/* makeHash generates MD5 hashes for given passwords */
func makeHash(le *logrus.Entry, cleartext string) string {
	hash := sha256.Sum256([]byte(cleartext))

	return hex.EncodeToString(hash[:])
}
