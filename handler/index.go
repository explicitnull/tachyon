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

	header := Header{
		Name: ctx.Value("username").(string),
	}

	// TODO: what is it?
	if ctx.Value("username").(string) == "dzhargalov" {
		header.Item10 = "disabled"
	}

	hdr, _ := template.ParseFiles("templates/hdr.htm")
	hdr.Execute(w, header)

	fmt.Fprintf(w, "<p>Welcome!</p>")

	ftr, _ := template.ParseFiles("templates/ftr.htm")
	ftr.Execute(w, nil)
}
