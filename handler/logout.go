package handler

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (g *Gateway) Logout(w http.ResponseWriter, r *http.Request) {
	deleteCookie(w)

	log.Info("user logged out")
	
	t, err := template.ParseFiles("templates/logout.htm")
	if err != nil {
		log.Errorf("%v", err)
	}
	t.Execute(w, nil)
}

func deleteCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "username",
		Value: "deleted",
		Path:  "/",
	}

	http.SetCookie(w, cookie)
}
