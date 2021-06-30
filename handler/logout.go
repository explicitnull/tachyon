package handler

import (
	"io"
	"net/http"
)

func (g *Gateway) Logout(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "logged out!\n")
}
