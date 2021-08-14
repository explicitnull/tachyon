package handler

import (
	"fmt"
	"net/http"
	"tacacs-webconsole/applogic"
	"tacacs-webconsole/types"
	"time"
)

func (g *Gateway) ShowAuthentication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	items, err := applogic.ShowAuthentication(le, g.aerospikeClient)
	if err != nil {
		le.WithError(err).Error("getting authentications failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	auts := types.Authentications{
		Items: items,
	}

	if len(items) == 0 {
		auts.NotFound = true
	} else {
		for range items {
			auts.Total++
		}
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	le.Debugf("%#v", auts)
	executeTemplate(le, w, "auth.htm", auts)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}

func (g *Gateway) SearchAuthentication(w http.ResponseWriter, r *http.Request) {
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

	items := applogic.SearchAuthentication(le, field, value, begin, end, g.aerospikeClient)

	auts := &types.Authentications{
		Items:     items,
		MoreItems: false,
	}

	if len(items) == 0 {
		auts.NotFound = true
	}

	if value != "" {
		auts.SearchValue = value
	} else if from != "" || to != "" {
		auts.SearchValue = fmt.Sprintf("%s - %s", from, to)
	}

	// counting summary
	for range items {
		auts.Total++
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "auth_search.htm", auts)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}
