package handler

import (
	"net/http"
	"tacacs-webconsole/applogic"
	"tacacs-webconsole/repository"
	"tacacs-webconsole/types"
)

const defaultAccountingRecordsPerPageLimit = 100

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

func (g *Gateway) SearchAccounting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	// if repository.GetRole(le, g.aerospikeClient, authenticatedUsername) != "superuser" ||
	// 	repository.GetRole(le, g.aerospikeClient, authenticatedUsername) != "manager" {
	// 	le.Warn("access forbidden")
	// 	http.Error(w, "access forbidden", http.StatusForbidden)
	// 	return
	// }

	field := r.PostFormValue("fld")
	value := r.PostFormValue("val")
	from := r.PostFormValue("from")
	to := r.PostFormValue("to")

	le.Debugf("field: %s, val: %s, from: %s, to: %s", field, value, from, to)

	if value == "" && from == "" && to == "" {
		le.Error("empty form")
		http.Error(w, "bad request: empty form", http.StatusBadRequest)
		return
	}

	items := applogic.SearchAccounting(le, field, value, from, to, g.aerospikeClient)

	acct := &types.AccountingRecords{
		Items:       items,
		MoreItems:   false,
		SearchValue: value,
	}

	if len(items) == 0 {
		acct.NotFound = true
	}

	// counting summary
	for range items {
		acct.Total++
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "acct_search.htm", acct)

	executeFooterTemplate(le, w)
}
