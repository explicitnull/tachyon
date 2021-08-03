package handler

import (
	"net/http"
	"tacasa-web/repository"
	"tacasa-web/types"
)

func (g *Gateway) ShowAuthentications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	items, err := repository.GetAuthentications(le, g.aerospikeClient)
	if err != nil {
		le.WithError(err).Error("getting authentications failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	auts := types.Authentications{
		Items: items,
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	le.Debugf("%#v", auts)
	executeTemplate(le, w, "auth.htm", auts)

	executeFooterTemplate(le, w)
}
