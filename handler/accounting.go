package handler

import (
	"net/http"
	"tacasa-web/repository"
	"tacasa-web/types"
)

func (g *Gateway) ShowAccounting(w http.ResponseWriter, r *http.Request) {
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

	items, _ := repository.GetAccounting(le, g.aerospikeClient)

	acct := &types.AccountingRecs{
		Items: items,
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "acct.htm", acct)

	executeFooterTemplate(le, w)
}
