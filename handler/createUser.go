package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"tachyon-web/repository"

	"github.com/aerospike/aerospike-client-go"
	"github.com/dchest/uniuri"
	"github.com/sirupsen/logrus"
)

type FormOptions struct {
	Opt1 []string
	Opt2 []string
}

type Req struct {
	Name   string
	Mail   string
	Subdiv string
	Permis string

	SubdivID int
	PermisID int
}

type CreatedUser struct {
	Name      string
	Cleartext string
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
	form.Opt1 = makeFormSelect("subdiv", "subdiv", "none")
	form.Opt2 = makeFormSelect("prm", "prm", "none")

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
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	// b, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	le.Error(err)
	// }
	// fmt.Println(string(b))

	r.ParseForm()
	req := Req{}
	req.Name = r.PostFormValue("username")
	req.Mail = r.PostFormValue("mail")
	if req.Mail == "" {
		req.Mail = "unknown"
	}
	req.Subdiv = r.PostFormValue("subdiv")
	req.Permis = r.PostFormValue("perm")

	le.Debugf("%#v", req)

	cleartext, err := createUserDo(le, g.aerospikeClient, req, authenticatedUsername)
	if err != nil {
		http.Error(w, "creating user failed", http.StatusInternalServerError)
		return
	}

	cu := CreatedUser{
		Name:      req.Name,
		Cleartext: cleartext,
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	mid, err := template.ParseFiles("templates/usercreated.htm")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mid.Execute(w, cu)

	executeFooterTemplate(le, w)
}

func genPass() string {
	return uniuri.NewLen(10)
}

func createUserDo(le *logrus.Entry, aClient *aerospike.Client, req Req, authenticatedUsername string) (string, error) {
	// normalization
	subdivID, err := repository.GetSubdivisionID(le, aClient, req.Subdiv)
	if err != nil {
		le.WithError(err).Error("getting subdivision ID failed")
		return "", err
	}

	permisID, err := repository.GetPermId(le, aClient, req.Permis)
	if err != nil {
		le.WithError(err).Error("getting permission ID failed")
		return "", err
	}

	cleartext := genPass()
	le.Debug(cleartext)
	hash := makeHash(le, cleartext)
	le.Debug(hash)

	err = repository.CreateUser(le, aClient, req.Name, hash, req.Mail, authenticatedUsername, permisID, subdivID)
	if err != nil {
		le.WithError(err).Errorf("error creating user")
		return "", err
	}

	le.WithField("username", req.Name).Info("user created")

	return cleartext, nil
}
