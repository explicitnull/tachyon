package handler

import (
	"fmt"
	"net/http"
	"tacacs-webconsole/repository"
	"tacacs-webconsole/types"

	"github.com/aerospike/aerospike-client-go"
	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const defaultAccountsPerPageLimit = 10

func (g *Gateway) ShowAccounts(w http.ResponseWriter, r *http.Request) {
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

	items, err := repository.GetAccounts(le, g.aerospikeClient)
	if err != nil {
		le.WithError(err).Error("getting accounts failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// preparing response
	accounts := &types.Accounts{
		Items:     items,
		MoreItems: false,
	}

	// counting summary
	for _, v := range items {
		accounts.Total++

		switch v.Status {
		case types.AccountStatusActive:
			accounts.Active++
		case types.AccountStatusPasswordNotChanged:
			accounts.PasswordNotChanged++
		case types.AccountStatusSuspended:
			accounts.Suspended++
		}

		// if v.LastSignedInTimestamp == "" {
		// 	accounts.NeverSignedIn++
		// }
	}

	if accounts.Total > defaultAccountsPerPageLimit {
		accounts.MoreItems = true
		accounts.ItemsPerPageLimit = defaultAccountsPerPageLimit
	}

	executeHeaderTemplate(le, w, username)

	executeTemplate(le, w, "accounts.htm", accounts)

	executeFooterTemplate(le, w)
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

	executeTemplate(le, w, "account_new.htm", form)

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
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// getting account data from DB
	acc, err := repository.GetAccountByName(le, g.aerospikeClient, name)
	if err != nil {
		le.WithError(err).Error("getting account failed")
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// filling template data
	acc.SubdivisionsList = repository.GetSubdivisionsList(le, g.aerospikeClient)
	acc.PermissionsList = repository.GetPermissionsList(le, g.aerospikeClient)

	// writing response
	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "account.htm", acc)

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
	name, ok := vars["name"]
	if !ok {
		le.Error(noIDinURL)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	fac := &types.Account{
		Name:        name,
		Cleartext:   r.PostFormValue("pwd"),
		Subdivision: r.PostFormValue("subdiv"),
		Permission:  r.PostFormValue("perm"),
		Mail:        r.PostFormValue("m"),
	}

	_, act := r.Form["active"]
	if act {
		fac.Status = "active"
	}

	// getting account data from DB
	dbac, err := repository.GetAccountByName(le, g.aerospikeClient, fac.Name)
	if err != nil {
		le.WithError(err).Error("getting account failed")
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// applying changes
	if fac.Cleartext != "" {
		hash := makeHash(le, fac.Cleartext)
		le.Debug(hash)

		err = repository.SetPassword(g.aerospikeClient, fac.Name, hash)
		if err != nil {
			http.Error(w, databaseError, http.StatusInternalServerError)
			return
		}
	}

	if fac.Subdivision != emptySelect && fac.Subdivision != dbac.Subdivision {
		err = repository.SetSubdivision(fac.Name, fac.Subdivision)
		if err != nil {
			http.Error(w, databaseError, http.StatusInternalServerError)
			return
		}
	}

	if fac.Permission != emptySelect && fac.Permission != dbac.Permission {
		err = repository.SetSubdivision(fac.Name, fac.Subdivision)
		if err != nil {
			http.Error(w, databaseError, http.StatusInternalServerError)
			return
		}
	}

	if fac.Mail != "" && fac.Mail != dbac.Mail {
		err = repository.SetMail(fac.Name, fac.Mail)
		if err != nil {
			http.Error(w, databaseError, http.StatusInternalServerError)
			return
		}
	}

	if fac.Status != dbac.Status {
		err = repository.SetAccountStatus(le, fac.Name, fac.Status)
		if err != nil {
			http.Error(w, databaseError, http.StatusInternalServerError)
		}
	}

	// redirecting back
	http.Redirect(w, r, r.URL.String()+"?from=editing", http.StatusTemporaryRedirect)
}

func (g *Gateway) DisableAccount(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// applying changes
	err := repository.SetAccountStatus(le, name, types.AccountStatusSuspended)
	if err != nil {
		le.WithError(err).Error("setting account status failure")
		http.Error(w, "setting account status failure", http.StatusInternalServerError)
		return
	}

	acc := types.Account{
		Name: name,
	}

	// writing response
	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "account_disabled.htm", acc)

	executeFooterTemplate(le, w)
}
