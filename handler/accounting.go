package handler

import (
	"net/http"
	"tacacs-webconsole/repository"
	"tacacs-webconsole/types"
)

const defaultAccountingRecordsPerPageLimit = 10

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

	items, err := repository.GetAccounting(le, g.aerospikeClient)
	if err != nil {
		le.WithError(err).Error("getting accounting failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	acct := &types.AccountingRecords{
		Items:     items,
		MoreItems: false,
	}

	// counting summary
	for _, _ = range items {
		acct.Total++
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "acct.htm", acct)

	executeFooterTemplate(le, w)
}
