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
			fmt.Fprintf(w, "error parsing template")
			le.Fatal(err)
		}
		t.Execute(w, nil)
		return
	}

	username, ok := ctx.Value("username").(string)
	if !ok {
		le.Warn("no username in context")
		fmt.Fprintf(w, "access forbidden")
		return
	}

	executeHeaderTemplate(le, w, username)

	fmt.Fprintf(w, "<p>Welcome!</p>")

	executeFooterTemplate(le, w)
}
