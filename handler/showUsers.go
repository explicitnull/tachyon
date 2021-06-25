package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type UserSummary struct {
	Total  int
	Active int
}

type User struct {
	Id         int
	Name       string
	Hash       string // Password hash
	Cleartext  string
	Subdiv     string
	Prm        string
	Mail       string
	Active     bool
	ActiveSt   string
	ActiveBox  string // Is "checked" or "" for HTML form
	CreaTime   string // Full time form
	CreaTimeS  string // Short time form
	CreaBy     string
	PassChd    bool
	SubdivList []string
	PrmList    []string
}

func (g *Gateway) ShowUsers(w http.ResponseWriter, r *http.Request) {
	sum := new(UserSummary)

	err := g.db.QueryRow("SELECT COUNT(*) FROM usr").Scan(&sum.Total)
	if err != nil {
		log.Errorf("quering total users count failed: %v", err)
	}

	err = g.db.QueryRow("SELECT COUNT(*) FROM usr WHERE act=true").Scan(&sum.Active)
	if err != nil {
		log.Errorf("quering active users count failed: %v", err)
	}

	rows, err := g.db.Query("SELECT u.uid, u.username, p.prm, u.mail, d.subdiv, u.created, u.act, u.pass_chd, u.created_by FROM usr u JOIN prm p ON (u.prm_id = p.prm_id) JOIN subdiv d ON (u.subdiv_id = d.l2id) ORDER BY p.prm, u.username")
	if err != nil {
		log.Errorf("quering all users failed: %v", err)
		fmt.Fprintf(w, "database query failed")
		return
	}
	defer rows.Close()

	/* Writing response */
	header := Header{
		Name: "furai", // FIXME
	}

	hdr, err := template.ParseFiles("templates/hdr.htm")
	if err != nil {
		log.Errorf("template parsing failed: %v", err)
	}
	hdr.Execute(w, header)

	main, err := template.ParseFiles("templates/users.htm")
	if err != nil {
		log.Errorf("template parsing failed: %v", err)
	}
	main.Execute(w, sum)

	var flags [2]string

	u := new(User)

	for rows.Next() {
		err = rows.Scan(&u.Id, &u.Name, &u.Prm, &u.Mail, &u.Subdiv, &u.CreaTime, &u.Active, &u.PassChd, &u.CreaBy)
		if err != nil {
			log.WithError(err).Error("scanning sql rows failed")
		}

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
			log.WithError(err).Error("parsing time failed")
		}

		u.CreaTimeS = t.Format(tShort)

		fmt.Fprintf(w, `<tr><td><a href="/edituser/%v/">%s</a></td>`, u.Id, u.Name)
		fmt.Fprintf(w, `<td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>`, u.Prm, u.Subdiv, u.Mail, u.CreaTimeS, flags[0], flags[1], u.CreaBy)
	}
	fmt.Fprintln(w, "</table></div>")

	ftr, err := template.ParseFiles("templates/ftr-to-top.htm")
	if err != nil {
		log.Errorf("template parsing failed: %v", err)
	}
	ftr.Execute(w, nil)
}
