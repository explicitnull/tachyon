package main

import (
	"net/http"
	"os"
	"tachyon-web/handler"
	"tachyon-web/middleware"
	"tachyon-web/options"

	"github.com/aerospike/aerospike-client-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
)

const appName = "tacacs-webconsole"

var (
	host      = "13.48.3.15"
	port      = 3000
	namespace = "tacacs"
)

func init() {
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	// TODO: add govvv
	log.Warnf("starting %s...", appName)

	appOptions := new(options.Options)

	err := options.Load(appOptions)
	if err != nil {
		log.Fatalf("loading options failed: %v", err)
	}

	// db := database.Open(appOptions.DbHost, appOptions.DbName, appOptions.DbName, appOptions.DbPassword)

	aerospikeClient, err := aerospike.NewClient(host, port)
	if err != nil {
		log.Fatalf("aerospike init failed: %v", err)
	}

	var hashKey = []byte("secret")
	var blockKey = []byte("1234567890123456")
	sCookie := cookieInit(hashKey, blockKey)

	// handlers
	g, err := handler.NewGateway(appOptions, aerospikeClient, sCookie)
	if err != nil {
		log.Fatalf("handler init failed: %v", err)
	}

	// middleware
	m, err := middleware.NewMiddleware(appOptions, sCookie)
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

	log.Warn("listening http on port 8000")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("starting HTTP server failed: %v", err)
	}
}

func cookieInit(hashKey, blockKey []byte) *securecookie.SecureCookie {
	return securecookie.New(hashKey, blockKey)
}
