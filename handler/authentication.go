package handler

import (
	"net/http"
	"tacacs-webconsole/applogic"
	"tacacs-webconsole/types"
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
