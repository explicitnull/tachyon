package handler

import (
	"net/http"
	"tachyon/repository"
	"tachyon/types"
)

func (g *Gateway) ShowLockouts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	items, _ := repository.GetLockouts(le, g.aerospikeClient)

	lockouts := &types.Lockouts{
		Items: items,
	}

	// counting summary
	for range items {
		lockouts.Total++
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "lockout.htm", lockouts)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}
