package handler

import (
	"net/http"
	"tachyon/repository"
	"tachyon/types"
)

func (g *Gateway) ShowEquipment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	// TODO: eqp
	items, _ := repository.GetSubdivisions(le, g.aerospikeClient)

	subdivs := &types.Subdivisions{
		Items: items,
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "eqp.htm", subdivs)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}
