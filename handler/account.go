package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"tachyon-web/repository"
	"tachyon-web/types"

	"github.com/aerospike/aerospike-client-go"
	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (g *Gateway) ShowUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	if repository.GetRole(le, g.aerospikeClient, username) == "none" {
		le.Warn("access forbidden")
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	sum := repository.GetUserCount(g.aerospikeClient)
	aa, _ := repository.GetUsers(le, g.aerospikeClient)

	executeHeaderTemplate(le, w, username)

	mid, err := template.ParseFiles("templates/users.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
		return
	}
	mid.Execute(w, sum)

	for _, a := range aa {
		id := a.Name
		fmt.Fprintf(w, `<tr><td><a href="/edituser/%v/">%s</a></td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>`,
			id, a.Name, a.Permission, a.Subdivision, a.Mail, a.CreatedTimestamp, a.Status, a.PasswordSetTimestamp, a.CreatedBy)
	}
	fmt.Fprintln(w, "</table></div>")

	ftr, err := template.ParseFiles("templates/ftr-to-top.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
	}
	ftr.Execute(w, nil)
}

type FormOptions struct {
	Opt1 []string
	Opt2 []string
}

type UserCreated struct {
	Name      string
	Cleartext string
}

func (g *Gateway) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	if repository.GetRole(le, g.aerospikeClient, username) != "admin" {
		le.Warn("access forbidden")
		fmt.Fprintf(w, "access forbidden")
		return
	}

	executeHeaderTemplate(le, w, username)

	form := new(FormOptions)
	form.Opt1 = repository.GetSubdivisionsList(le, g.aerospikeClient)
	form.Opt2 = repository.GetPermissionsList(le, g.aerospikeClient)

	executeTemplate(le, w, "newuser.htm", form)

	executeFooterTemplate(le, w)
}

func (g *Gateway) CreateUserAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
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
	acc := types.Account{
		Name:        r.PostFormValue("username"),
		Mail:        r.PostFormValue("mail"),
		Subdivision: r.PostFormValue("subdiv"),
		Permission:  r.PostFormValue("perm"),
	}

	if acc.Mail == "" {
		acc.Mail = "unknown"
	}

	le.Debugf("%#v", acc)

	cleartext, err := createUserAction(le, g.aerospikeClient, acc, authenticatedUsername)
	if err != nil {
		http.Error(w, "creating user failed", http.StatusInternalServerError)
		return
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	uc := UserCreated{
		Name:      acc.Name,
		Cleartext: cleartext,
	}
	executeTemplate(le, w, "usercreated.htm", uc)

	executeFooterTemplate(le, w)
}

func genPass() string {
	return uniuri.NewLen(10)
}

func createUserAction(le *logrus.Entry, aClient *aerospike.Client, req types.Account, authenticatedUsername string) (string, error) {
	// normalization
	subdivID, err := repository.GetSubdivisionID(le, aClient, req.Subdivision)
	if err != nil {
		le.WithError(err).Error("getting subdivision ID failed")
		return "", err
	}

	permisID, err := repository.GetPermId(le, aClient, req.Permission)
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

func (g *Gateway) EditAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	if repository.GetRole(le, g.aerospikeClient, authenticatedUsername) != "admin" {
		le.Warn("access forbidden")
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	// parsing request
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		le.Error(noIDinURL)
	}

	// getting account data from DB
	acc, err := repository.GetAccountByName(le, g.aerospikeClient, name)
	if err != nil {
		le.WithError(err).Error("getting account failed")
		http.Error(w, "access forbidden", http.StatusForbidden)
	}

	// filling template data
	td := types.AccountTemplateData{
		Name:   acc.Name,
		Mail:   acc.Mail,
		Subdiv: acc.Subdivision,
		Prm:    acc.Permission,
	}

	if acc.Status == "active" {
		td.ActiveBox = "checked"
	}

	td.SubdivList = repository.GetSubdivisionsList(le, g.aerospikeClient)
	td.PrmList = repository.GetPermissionsList(le, g.aerospikeClient)

	// writing response
	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "user.htm", td)

	executeFooterTemplate(le, w)
}

func (g *Gateway) EditAccountAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	if repository.GetRole(le, g.aerospikeClient, authenticatedUsername) != "admin" {
		le.Warn("access forbidden")
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	// parsing request
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		le.Error(noIDinURL)
	}

	r.ParseForm()
	fac := &types.Account{
		Password:    r.PostFormValue("pwd"),
		Subdivision: r.PostFormValue("subdiv"),
		Permission:  r.PostFormValue("perm"),
		Mail:        r.PostFormValue("m"),
	}

	_, act := r.Form["active"]
	if act {
		fac.Status = "active"
	}

	// getting account data from DB
	dbac, err := repository.GetAccountByName(le, g.aerospikeClient, username)
	if err != nil {
		le.WithError(err).Error("getting account failed")
		http.Error(w, "access forbidden", http.StatusForbidden)
	}

	// applying changes
	if fac.Password != "" {
		hash := makeHash(le, fac.Password)
		le.Debug(hash)

		err = repository.SetPassword(g.aerospikeClient, username, hash)
		if err != nil {
			http.Error(w, databaseError, http.StatusInternalServerError)
		}
	}

	if fac.Subdivision != emptySelect && fac.Subdivision != dbac.Subdivision {
		err = repository.SetSubdivision(username, fac.Subdivision)
		if err != nil {
			http.Error(w, databaseError, http.StatusInternalServerError)
		}
	}

	if fac.Permission != emptySelect && fac.Permission != dbac.Permission {
		err = repository.SetSubdivision(username, fac.Subdivision)
		if err != nil {
			http.Error(w, databaseError, http.StatusInternalServerError)
		}
	}

	if fac.Mail != "" && fac.Mail != dbac.Mail {
		err = repository.SetMail(username, fac.Mail)
		if err != nil {
			http.Error(w, databaseError, http.StatusInternalServerError)
		}
	}

	if fac.Status != dbac.Status {
		err = repository.SetAccountStatus(le, username, fac.Status)
		if err != nil {
			http.Error(w, databaseError, http.StatusInternalServerError)
		}
	}

	// redirecting back
	http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
}
