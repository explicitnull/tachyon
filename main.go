package main

import (
	"net/http"
	"os"
	"tachyon-web/database"
	"tachyon-web/handler"
	"tachyon-web/middleware"
	"tachyon-web/options"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const appName = "tachyon-web"

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

	db := database.Open(appOptions.DbHost, appOptions.DbName, appOptions.DbName, appOptions.DbPassword)

	g, err := handler.NewGateway(appOptions, db)
	if err != nil {
		log.Fatalf("handler init failed: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", g.AppInfo)
	r.HandleFunc("/login/", g.Login).Methods("GET")
	r.HandleFunc("/login/", g.LoginDo).Methods("POST")
	r.HandleFunc("/logout/", g.Logout)
	r.HandleFunc("/users/", g.ShowUsers)

	r.Use(middleware.CheckCookie)

	// TODO: move serving assets to standalone proxy like nginx
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/", r)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("starting HTTP server failed: %v", err)
	}
}
