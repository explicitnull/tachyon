package handler

import (
	"fmt"
	"net/http"
	"tacacs-webconsole/applogic"
	"tacacs-webconsole/repository"
	"tacacs-webconsole/types"
	"time"
)

const defaultAccountingRecordsPerPageLimit = 100
const acctOffset = 10

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

	now := time.Now()
	begin := now.Add(-acctOffset * time.Minute)
	end := now

	le.Debugf("handler begin: %s, end: %s", begin, end)

	items, err := repository.GetAccountingWithTimeFilter(le, g.aerospikeClient, begin, end)
	if err != nil {
		le.WithError(err).Error("getting accounting failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	acct := &types.AccountingRecords{
		Items:     items,
		MoreItems: false,
	}

	if len(items) == 0 {
		acct.NotFound = true
	} else {
		for range items {
			acct.Total++
		}
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

	var (
		begin, end time.Time
		err        error
	)

	if value == "" {
		begin, err = time.Parse(types.TimeFormatSeconds, from)
		if err != nil {
			le.Error("wrong time format")
			http.Error(w, "bad request: wrong time format", http.StatusBadRequest)
			return
		}

		end, err = time.Parse(types.TimeFormatSeconds, to)
		if err != nil {
			le.Error("wrong time format")
			http.Error(w, "bad request: wrong time format", http.StatusBadRequest)
			return
		}
	}

	items := applogic.SearchAccounting(le, field, value, begin, end, g.aerospikeClient)

	acct := &types.AccountingRecords{
		Items:     items,
		MoreItems: false,
	}

	if len(items) == 0 {
		acct.NotFound = true
	}

	if value != "" {
		acct.SearchValue = value
	} else if from != "" && to != "" {
		acct.SearchValue = fmt.Sprintf("%s - %s", from, to)
	}

	// counting summary
	for range items {
		acct.Total++
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "acct_search.htm", acct)

	executeFooterTemplate(le, w)
}
