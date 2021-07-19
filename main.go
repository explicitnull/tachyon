package main

import (
	"net/http"
	"os"
	"tachyon-web/database"
	"tachyon-web/handler"
	"tachyon-web/middleware"
	"tachyon-web/options"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
)

const appName = "tachyon-web"

func init() {
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	// TODO: add govvv
	log.Warnf("Starting %s...", appName)

	appOptions := new(options.Options)

	err := options.Load(appOptions)
	if err != nil {
		log.Fatalf("loading options failed: %v", err)
	}

	db := database.Open(appOptions.DbHost, appOptions.DbName, appOptions.DbName, appOptions.DbPassword)

	var hashKey = []byte("secret")
	var blockKey = []byte("1234567890123456")
	sc := cookieInit(hashKey, blockKey)

	// handlers
	g, err := handler.NewGateway(appOptions, db, sc)
	if err != nil {
		log.Fatalf("handler init failed: %v", err)
	}

	// middleware
	m, err := middleware.NewMiddleware(appOptions, sc)
	if err != nil {
		log.Fatalf("middleware init failed: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", g.Index)
	r.HandleFunc("/login/", g.Login).Methods("GET")
	r.HandleFunc("/login/", g.LoginDo).Methods("POST")
	r.HandleFunc("/logout/", g.Logout)
	r.HandleFunc("/passwd/", g.ChangePassword).Methods("GET")
	r.HandleFunc("/passwd/", g.ChangePasswordDo).Methods("POST")
	r.HandleFunc("/newuser/", g.CreateUser).Methods("GET")
	r.HandleFunc("/newuser/", g.CreateUserDo).Methods("POST")

	r.HandleFunc("/users/", g.ShowUsers)

	r.Use(m.Log)
	r.Use(m.CheckCookie)

	// TODO: move serving assets to standalone proxy like nginx
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/", r)

	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("starting HTTP server failed: %v", err)
	}
}

func cookieInit(hashKey, blockKey []byte) *securecookie.SecureCookie {
	return securecookie.New(hashKey, blockKey)
}
