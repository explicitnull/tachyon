package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"tachyon/applogic"
)

func (g *Gateway) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	executeHeaderTemplate(le, w, username)

	mid, _ := template.ParseFiles("templates/chpass.htm")
	mid.Execute(w, nil)

	executeFooterTemplate(le, w)

	le.WithField("origin", "ChangePassword").Info("handled ok")
}

func (g *Gateway) ChangePasswordAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	r.ParseForm()
	f1 := r.Form["pass1"]
	f2 := r.Form["pass2"]
	pass := f1[0]
	passConfirm := f2[0]

	err := applogic.ChangePasswordAction(le, g.aerospikeClient, authenticatedUsername, pass, passConfirm, g.Options)
	if err == applogic.ErrPasswordsMismatch {
		fmt.Fprintf(w, "<p>Ошибка! Введенные пароли не совпадают.</p>")
		return
	} else if err == applogic.ErrBadCharacters {
		fmt.Fprintln(w, "<p>Ошибка! Пароль содержит недопустимые символы. Забыли переключить раскладку?</p>")
		return
	} else if err == applogic.ErrTooShort {
		fmt.Fprintln(w, "Ошибка! Пароль содержит слишком мало символов")
		return
	}

	fmt.Fprintf(w, "<p>Пароль изменен.</p>")

	le.WithField("origin", "ChangePasswordAction").Info("handled ok")
}
