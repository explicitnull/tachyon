package handler

import (
	"io"
	"net/http"
)

func (g *Gateway) AppInfo(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "this is tachyon-web\n")
}

func (g *Gateway) Login(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "logged in!\n")
}

func (g *Gateway) Logout(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "logged out!\n")
}
