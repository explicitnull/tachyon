package handler

import (
	"net/http"
	"tacacs-webconsole/repository"
	"tacacs-webconsole/types"
)

func (g *Gateway) ShowLockouts(w http.ResponseWriter, r *http.Request) {
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

	items, _ := repository.GetLockouts(le, g.aerospikeClient)

	lockouts := &types.Lockouts{
		Items: items,
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "lockout.htm", lockouts)

	executeFooterTemplate(le, w)
}
