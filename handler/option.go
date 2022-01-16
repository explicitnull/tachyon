package handler

import (
	"net/http"
	"tachyon/repository"
	"tachyon/types"
)

func (g *Gateway) ShowOptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	authenticatedUsername, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	o, err := repository.GetOptions(le, g.aerospikeClient)
	if err != nil {
		le.WithError(err).Error("getting options failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	t, err := repository.GetTokens(le, g.aerospikeClient)
	if err != nil {
		le.WithError(err).Error("getting tokens failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	opts := &types.Options{
		OptionItems: o,
		TokenItems:  t,
	}

	executeHeaderTemplate(le, w, authenticatedUsername)

	executeTemplate(le, w, "options.htm", opts)

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}

// func (g *Gateway) EditOptionsAction(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	le := getLogger(r)

// 	authenticatedUsername, ok := ctx.Value("username").(string)
// 	if !ok {
// 		le.Warn("no username in context")
// 		return
// 	}

// 	r.ParseForm()
// 	m := r.PostFormValue("minPassLen")
// 	minPassLen, err := strconv.Atoi(m)
// 	if err != nil {
// 		le.WithError(err).Error("parsing options failed")
// 		http.Error(w, badRequest, http.StatusBadRequest)
// 		return
// 	}

// 	fo := &types.Option{
// 		MinimumPasswordLength: minPassLen,
// 	}

// 	// getting account data from DB
// 	dbo, err := repository.GetOptions(le, g.aerospikeClient)
// 	if err != nil {
// 		le.WithError(err).Error("getting option failed")
// 		http.Error(w, serverError, http.StatusInternalServerError)
// 		return
// 	}

// 	// applying changes
// 	if fo.MinimumPasswordLength != dbo.MinimumPasswordLength {
// 		err = repository.SetOptionMinimumPasswordLength(le, fo.MinimumPasswordLength)
// 		if err != nil {
// 			http.Error(w, serverError, http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	// redirecting back
// 	http.Redirect(w, r, r.URL.String()+"?from=editing", http.StatusTemporaryRedirect)
// }
