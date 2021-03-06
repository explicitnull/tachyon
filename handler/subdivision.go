package handler

import (
	"net/http"
	"tachyon/repository"
	"tachyon/types"
)

func (g *Gateway) ShowSubdivisions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	items, _ := repository.GetSubdivisions(le, g.aerospikeClient)

	subdivs := &types.Subdivisions{
		Items: items,
	}

	// counting summary
	for _, v := range items {
		subdivs.Total++

		switch v.Status {
		case types.SubdivisionStatusActive:
			subdivs.Active++
		case types.SubdivisionStatusInactive:
			subdivs.Inactive++
		}
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "subdiv.htm", subdivs)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}

func (g *Gateway) CreateSubdivision(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	executeHeaderTemplate(le, w, username)

	executeTemplate(le, w, "subdiv_new.htm", nil)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}
