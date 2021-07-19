package handler

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"tachyon-web/repository"
)

type FormOptions struct {
	Opt1 []string
	Opt2 []string
}

func (g *Gateway) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		fmt.Fprintf(w, "access forbidden")
		return
	}

	if repository.GetRole(le, g.aerospikeClient, username) != "admin" {
		le.Warn("access forbidden")
		fmt.Fprintf(w, "access forbidden")
		return
	}

	executeHeaderTemplate(le, w, username)

	form := new(FormOptions)
	// form.Opt1 = makeFormSelect("subdiv", "subdiv", "none")
	// form.Opt2 = makeFormSelect("prm", "prm", "none")

	mid, err := template.ParseFiles("templates/newuser.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
		return
	}

	err = mid.Execute(w, form)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Template not found: %v", err)
		log.Printf("Template not found: %v", err)
	}

	executeFooterTemplate(le, w)
}

func (g *Gateway) CreateUserDo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		fmt.Fprintf(w, "access forbidden")
		return
	}

	if repository.GetRole(le, g.aerospikeClient, authenticatedUsername) != "admin" {
		le.Warn("access forbidden")
		fmt.Fprintf(w, "access forbidden")
		return
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	r.ParseForm()
	username := r.PostFormValue("username")
	mail := r.PostFormValue("mail")
	subdiv := r.PostFormValue("subdiv")
	prm := r.PostFormValue("perm")

	if mail == "" {
		mail = "unknown"
	}

	// normalization
	subdiv_id := repository.GetSubdivisionID(le, g.aerospikeClient, subdiv)
	prm_id, err := repository.GetPermId(le, g.aerospikeClient, prm)
	if err != nil {
		le.WithError(err).Error("getting permissions failed")
		return
	}

	hash := makeHash(le, genPass())

	err = repository.CreateUser(le, g.aerospikeClient, username, hash, mail, authenticatedUsername, prm_id, subdiv_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Template error for user %s: %s", authenticatedUsername, err)
		return
	}

	mid, err := template.ParseFiles("templates/usercreated.htm")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mid.Execute(w, nil)

	executeFooterTemplate(le, w)

	le.WithField("username", username).Info("user created")
}

func genPass() string {
	cmd := exec.Command("openssl", "rand", "-base64", "7")
	stdout, err := cmd.StdoutPipe()
	checkErr(err)

	cmd.Start()
	pipe := bufio.NewReader(stdout)

	line, _, err := pipe.ReadLine()
	checkErr(err)
	res := string(line)
	pass := strings.Replace(res, "=", "", 2)
	return pass
}
