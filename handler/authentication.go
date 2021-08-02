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

	items, _ := repository.GetAuthentications(le, g.aerospikeClient)

	auth := &types.Authentications{
		Items: items,
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "auth.htm", auth)

	executeFooterTemplate(le, w)
}
