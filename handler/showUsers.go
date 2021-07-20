package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"tachyon-web/repository"
)

func (g *Gateway) ShowUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	if repository.GetRole(le, g.aerospikeClient, username) == "none" {
		le.Warn("access forbidden")
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	sum := repository.GetUserCount(g.aerospikeClient)
	users, _ := repository.GetUsers(le, g.aerospikeClient)

	executeHeaderTemplate(le, w, username)

	mid, err := template.ParseFiles("templates/users.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
		return
	}
	mid.Execute(w, sum)

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

		// t, err := time.Parse(time.RFC3339Nano, u.CreaTime)
		// if err != nil {
		// 	le.WithError(err).Error("time parsing failed")
		// 	return
		// }

		// u.CreaTimeS = t.Format(timeShort)

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
