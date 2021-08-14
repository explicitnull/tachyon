package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

func (g *Gateway) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	if g.Options.Maintenance == "yes" {
		t, err := template.ParseFiles("templates/mntn.htm")
		if err != nil {
			le.WithError(err).Error("template parsing failed")
			http.Error(w, "template parsing failed", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
		return
	}

	username, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		return
	}

	executeHeaderTemplate(le, w, username)

	fmt.Fprintf(w, "<p>Welcome!</p>")

	executeFooterTemplate(le, w)

	le.Info("handled ok")
}
