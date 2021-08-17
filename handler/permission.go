package handler

import (
	"fmt"
	"net/http"
	"tacacs-webconsole/applogic"
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

	if repository.GetRole(le, g.aerospikeClient, authenticatedUsername) != "superuser" {
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

	le.Info("handled ok")
}

func (g *Gateway) CreatePermission(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	if repository.GetRole(le, g.aerospikeClient, username) != "superuser" {
		le.Warn("access forbidden")
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	executeHeaderTemplate(le, w, username)

	executeTemplate(le, w, "permission_new.htm", nil)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}

func (g *Gateway) CreatePermissionAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	if repository.GetRole(le, g.aerospikeClient, authenticatedUsername) != "superuser" {
		le.Warn("access forbidden")
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	r.ParseForm()
	perm := &types.Permission{
		Name:        r.PostFormValue("perm"),
		Description: r.PostFormValue("descr"),
		CreatedBy:   authenticatedUsername,
	}

	err := repository.CreatePermission(le, g.aerospikeClient, perm)
	if err != nil {
		http.Error(w, "creating user failed", http.StatusInternalServerError)
		return
	}

	notice := Notice{
		Title:   "New permission",
		Message: fmt.Sprintf("Permission \"%s\" created.", perm.Name),
	}

	// writing response
	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "notice.htm", notice)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}

func (g *Gateway) EditPermission(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	if repository.GetRole(le, g.aerospikeClient, authenticatedUsername) != "superuser" {
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

	// getting data from DB
	perm, err := repository.GetPermissionByName(le, g.aerospikeClient, name)
	if err != nil {
		le.WithError(err).Error("getting permission failed")
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// writing response
	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "permission.htm", perm)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}

func (g *Gateway) EditPermissionAction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	if repository.GetRole(le, g.aerospikeClient, authenticatedUsername) != "superuser" {
		le.Warn("access forbidden")
		http.Error(w, accessForbidden, http.StatusForbidden)
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

	r.ParseForm()
	fperm := types.Permission{
		Name:        name,
		Description: r.PostFormValue("descr"),
		Status:      r.PostFormValue("status"),
	}

	// _, act := r.Form["active"]
	// if act {
	// 	fperm.Status = "active"
	// }

	// applying changes
	err := applogic.EditPermissionAction(le, g.aerospikeClient, fperm)
	if err != nil {
		le.WithError(err).Error("setting permission properties failure")
		http.Error(w, serverError, http.StatusInternalServerError)
		return
	}

	notice := Notice{
		Title:   fmt.Sprintf("Permission \"%s\"", name),
		Message: "Changes saved",
	}

	// writing response
	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "notice.htm", notice)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}
