package handler

import (
	"net/http"
	"tacacs-webconsole/repository"
	"tacacs-webconsole/types"

	"github.com/gorilla/mux"
)

type PermissionCreated struct {
	Name string
}

func (g *Gateway) ShowPermissions(w http.ResponseWriter, r *http.Request) {
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

	items, _ := repository.GetPermissions(le, g.aerospikeClient)

	perms := &types.Permissions{
		Items: items,
	}

	// counting summary
	for _, v := range items {
		perms.Total++

		switch v.Status {
		case types.PermissionStatusActive:
			perms.Active++
		case types.PermissionStatusInactive:
			perms.Inactive++
		}
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "permissions.htm", perms)

	executeFooterTemplate(le, w)
}

func (g *Gateway) CreatePermission(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	if repository.GetRole(le, g.aerospikeClient, username) != "admin" {
		le.Warn("access forbidden")
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	executeHeaderTemplate(le, w, username)

	executeTemplate(le, w, "newpermission.htm", nil)

	executeFooterTemplate(le, w)
}

func (g *Gateway) CreatePermissionAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	if repository.GetRole(le, g.aerospikeClient, username) != "admin" {
		le.Warn("access forbidden")
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	r.ParseForm()
	perm := &types.Permission{
		Name:        r.PostFormValue("perm"),
		Description: r.PostFormValue("descr"),
		CreatedBy:   username,
	}

	err := repository.CreatePermission(le, g.aerospikeClient, perm)
	if err != nil {
		http.Error(w, "creating user failed", http.StatusInternalServerError)
		return
	}

	pc := PermissionCreated{
		Name: perm.Name,
	}
	executeTemplate(le, w, "permissioncreated.htm", pc)

	executeFooterTemplate(le, w)
}

func (g *Gateway) EditPermission(w http.ResponseWriter, r *http.Request) {
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

	// getting data from DB
	perm, err := repository.GetPermissionByName(le, g.aerospikeClient, name)
	if err != nil {
		le.WithError(err).Error("getting permission failed")
		http.Error(w, "access forbidden", http.StatusForbidden)
	}

	// writing response
	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "permission.htm", perm)

	executeFooterTemplate(le, w)
}

func (g *Gateway) EditPermissionAction(w http.ResponseWriter, r *http.Request) {
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

	r.ParseForm()
	fperm := &types.Account{
		Name:        vars["name"],
		Cleartext:   r.PostFormValue("descr"),
		Subdivision: r.PostFormValue("status"),
	}

	_, act := r.Form["active"]
	if act {
		fperm.Status = "active"
	}

	// getting account data from DB
	dbac, err := repository.GetAccountByName(le, g.aerospikeClient, fperm.Name)
	if err != nil {
		le.WithError(err).Error("getting account failed")
		http.Error(w, "access forbidden", http.StatusForbidden)
	}

	// applying changes
	if fperm.Mail != "" && fperm.Mail != dbac.Mail {
		err = repository.SetMail(fperm.Name, fperm.Mail)
		if err != nil {
			http.Error(w, databaseError, http.StatusInternalServerError)
		}
	}

	if fperm.Status != dbac.Status {
		err = repository.SetAccountStatus(le, fperm.Name, fperm.Status)
		if err != nil {
			http.Error(w, databaseError, http.StatusInternalServerError)
		}
	}

	// redirecting back
	http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
}
