package handler

import (
	"html/template"
	"net/http"
)

func (g *Gateway) Logout(w http.ResponseWriter, r *http.Request) {
	le := getLogger(r)

	deleteCookie(w)

	le.Info("user logged out")

	t, err := template.ParseFiles("templates/logout.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
		http.Error(w, "template parsing failed", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)

	le.Info("handled ok")
}

func deleteCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "username",
		Value: "deleted",
		Path:  "/",
	}

	http.SetCookie(w, cookie)
}
