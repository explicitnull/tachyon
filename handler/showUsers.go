package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"tachyon-web/repository"
	"time"
)

func (g *Gateway) ShowUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username := ctx.Value("username").(string)

	if repository.CheckRole(le, g.db, username) == "none" {
		le.Warn("access forbidden")
		fmt.Fprintf(w, "access forbidden")
		return
	}

	sum := repository.GetUserCount(g.db)

	header := Header{
		Name: ctx.Value("username").(string),
	}

	hdr, err := template.ParseFiles("templates/hdr.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
	}
	hdr.Execute(w, header)

	main, err := template.ParseFiles("templates/users.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
	}
	main.Execute(w, sum)

	users := repository.GetUsers(g.db)

	var flags [2]string

	for _, u := range users {
		if u.Active {
			flags[0] = "Активна"
		} else {
			flags[0] = "Неактивна"
		}

		if u.PassChd {
			flags[1] = "Постоянный"
		} else {
			flags[1] = "Временный"
		}

		t, err := time.Parse(time.RFC3339Nano, u.CreaTime)
		if err != nil {
			le.WithError(err).Error("template parsing failed")
		}

		u.CreaTimeS = t.Format(timeShort)

		fmt.Fprintf(w, `<tr><td><a href="/edituser/%v/">%s</a></td>`, u.Id, u.Name)
		fmt.Fprintf(w, `<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>`, u.Prm, u.Subdiv, u.Mail, u.CreaTimeS, flags[0], flags[1], u.CreaBy)
	}
	fmt.Fprintln(w, "</table></div>")

	ftr, err := template.ParseFiles("templates/ftr-to-top.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
	}
	ftr.Execute(w, nil)
}
