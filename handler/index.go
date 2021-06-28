package handler

import (
	"fmt"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (g *Gateway) Index(w http.ResponseWriter, r *http.Request) {
	if g.Options.Maintenance == "yes" {
		t, err := template.ParseFiles("templates/mntn.htm")
		if err != nil {
			fmt.Fprintf(w, "error parsing template")
			log.Fatal(err)
		}
		t.Execute(w, nil)
		return
	}

	header := Header{
		Name: "furai", // FIXME:
	}

	// TODO: what is it?
	// if authUser == "dzhargalov" {
	// 	header.Item10 = "disabled"
	// }

	hdr, _ := template.ParseFiles("templates/hdr.htm")
	hdr.Execute(w, header)

	fmt.Fprintf(w, "<p>Welcome!</p>")

	ftr, _ := template.ParseFiles("templates/ftr.htm")
	ftr.Execute(w, nil)
}
