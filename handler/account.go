package handler

import (
	"fmt"
	"net/http"
	"tachyon/applogic"
	"tachyon/repository"
	"tachyon/types"

	"github.com/gorilla/mux"
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

	items, err := repository.GetAccounts(ctx, g.aerospikeClient)
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

	le.WithField("origin", "ShowAccounts").Info("handled ok")
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

	executeHeaderTemplate(le, w, username)

	form := new(FormOptions)
	form.Opt1 = repository.GetSubdivisionsList(le, g.aerospikeClient)
	form.Opt2 = repository.GetPermissionsList(le, g.aerospikeClient)

	executeTemplate(le, w, "account_new.htm", form)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}

func (g *Gateway) CreateUserAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	r.ParseForm()
	acc := types.Account{
		Name:        r.PostFormValue("username"),
		Subdivision: r.PostFormValue("subdiv"),
		Permission:  r.PostFormValue("perm"),
		Mail:        r.PostFormValue("mail"),
		Status:      r.PostFormValue("status"),
	}

	if acc.Mail == "" {
		acc.Mail = "unknown"
	}

	le.Debugf("%#v", acc)

	cleartext, err := applogic.CreateUserAction(le, g.aerospikeClient, acc, authenticatedUsername)
	if err != nil {
		http.Error(w, "creating user failed", http.StatusInternalServerError)
		return
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	notice := Notice{
		Title:   "New account",
		Message: fmt.Sprintf("User \"%s\" created. Password: \"%s\". Network access will be provided in a few minutes.", acc.Name, cleartext),
	}

	executeTemplate(le, w, "notice.htm", notice)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}

func (g *Gateway) EditAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	// parsing request
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		le.Error(noIDinURL)
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	// getting account data from DB
	acc, err := repository.GetAccountByName(ctx, g.aerospikeClient, name)
	if err != nil {
		le.WithError(err).Error("getting account failed")
		http.Error(w, serverError, http.StatusInternalServerError)
		return
	}

	// filling template data
	acc.SubdivisionsList = repository.GetSubdivisionsList(le, g.aerospikeClient)
	acc.PermissionsList = repository.GetPermissionsList(le, g.aerospikeClient)

	// writing response
	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "account.htm", acc)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}

func (g *Gateway) EditAccountAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
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
		Status:      r.PostFormValue("status"),
		UILevel:     r.PostFormValue("ui_level"),
	}

	err := applogic.EditAccountAction(le, g.aerospikeClient, fac)
	if err != nil {
		http.Error(w, serverError, http.StatusInternalServerError)
		return
	}

	notice := Notice{
		Title:   fmt.Sprintf("Account \"%s\"", name),
		Message: "Changes saved",
	}

	// writing response
	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "notice.htm", notice)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}

func (g *Gateway) RemoveAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
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
	err := applogic.RemoveAccount(le, g.aerospikeClient, name, authenticatedUsername)
	if err != nil {
		le.WithError(err).Error("setting account status failure")
		http.Error(w, "setting account status failure", http.StatusInternalServerError)
		return
	}

	notice := Notice{
		Title:   fmt.Sprintf("Account \"%s\"", name),
		Message: "Account removed",
	}

	// writing response
	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "notice.htm", notice)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}
