package handler

import (
	"fmt"
	"net/http"
	"tachyon/applogic"
	"tachyon/types"
	"time"
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

	items, err := applogic.ShowAccounting(le, g.aerospikeClient)
	if err != nil {
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

	le.Info("handled ok")
}

func (g *Gateway) SearchAccounting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

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

	le.Info("handled ok")
}
