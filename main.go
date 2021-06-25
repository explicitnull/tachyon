package main

import (
	"net/http"
	"os"
	"tac-gateway/handler"
	"tac-gateway/options"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	log "github.com/sirupsen/logrus"
)

const appName = "tac-gateway"

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Warnf("Starting %s...", appName)

	appOptions := new(options.Options)

	err := options.Load(appOptions)
	if err != nil {
		log.Fatalf("loading options failed: %v", err)
	}

	db := db.dbconf()

	g, err := handler.NewGateway(appOptions, db)
	if err != nil {
		log.Fatalf("handler init failed: %v", err)
	}

	mx := mux.NewRouter()
	mx.HandleFunc("/", g.AppInfo)
	mx.HandleFunc("/login", g.Login)
	mx.HandleFunc("/logout", g.Logout)
	mx.HandleFunc("/users/", alice.New(g.CheckExtendedAccess).Then(g.ShowUsers))

	http.Handle("/", mx)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("starting HTTP server failed: %v", err)
	}
}
