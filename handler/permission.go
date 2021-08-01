package handler

import (
	"net/http"
	"tacasa-web/repository"
	"tacasa-web/types"
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

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "prm.htm", perms)

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
