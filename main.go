package main

import (
	"net/http"
	"os"
	"tacacs-webconsole/handler"
	"tacacs-webconsole/middleware"
	"tacacs-webconsole/options"

	"github.com/aerospike/aerospike-client-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"
)

const appName = "tacacs-webconsole"

// TODO: move to env
var (
	host = "13.48.3.15"
	port = 3000
)

func main() {
	loggingInit()

	// TODO: add govvv
	log.Warnf("starting %s", appName)
	log.Warnf("establishing connection to database")

	// TODO: pg support
	// db := database.Open(appOptions.DbHost, appOptions.DbName, appOptions.DbName, appOptions.DbPassword)

	aerospikeClient, err := aerospike.NewClient(host, port)
	if err != nil {
		log.Fatalf("aerospike init failed: %v", err)
	}

	log.Warnf("receiving settings")

	appOptions := new(options.Options)
	appOptions.MinPassLen = 9

	err = options.Load(appOptions)
	if err != nil {
		log.Fatalf("loading options failed: %v", err)
	}

	// TODO: move to config
	hashKey := []byte("secret")
	blockKey := []byte("1234567890123456")
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
	r.HandleFunc("/", g.Index).Methods("GET")
	r.HandleFunc("/login/", g.Login).Methods("GET")
	r.HandleFunc("/login/", g.LoginAction).Methods("POST")
	r.HandleFunc("/logout/", g.Logout).Methods("GET")

	// tacplus configuration management handlers

	r.HandleFunc("/myaccount/", g.ChangePassword).Methods("GET")
	r.HandleFunc("/myaccount/", g.ChangePasswordAction).Methods("POST")

	r.HandleFunc("/users/", g.ShowAccounts).Methods("GET")
	r.HandleFunc("/newuser/", g.CreateUser).Methods("GET")
	r.HandleFunc("/newuser/", g.CreateUserAction).Methods("POST")
	r.HandleFunc("/edituser/{name}/", g.EditAccount).Methods("GET")
	r.HandleFunc("/edituser/{name}/", g.EditAccountAction).Methods("POST")
	r.HandleFunc("/disableaccount/{name}/", g.DisableAccount)

	r.HandleFunc("/permissions/", g.ShowPermissions).Methods("GET")
	r.HandleFunc("/newpermission/", g.CreatePermission).Methods("GET")
	r.HandleFunc("/newpermission/", g.CreatePermissionAction).Methods("POST")
	r.HandleFunc("/editpermission/{name}/", g.EditPermission).Methods("GET")
	r.HandleFunc("/editpermission/{name}/", g.EditPermissionAction).Methods("POST")

	r.HandleFunc("/subdiv/", g.ShowSubdivisions).Methods("GET")
	r.HandleFunc("/subdiv/new/", g.CreateSubdivision).Methods("GET")

	r.HandleFunc("/equipment/", g.ShowEquipment).Methods("GET")

	// tacplus logs handlers

	r.HandleFunc("/auth/", g.ShowAuthentication).Methods("GET")
	r.HandleFunc("/auth_search/", g.SearchAuthentication).Methods("POST")

	r.HandleFunc("/acct/", g.ShowAccounting).Methods("GET")
	r.HandleFunc("/acct-search/", g.SearchAccounting).Methods("POST")

	// antibruteforce handler
	r.HandleFunc("/lockout/", g.ShowLockouts)

	r.HandleFunc("/settings/", g.ShowOptions).Methods("GET")

	r.Use(m.Log)
	r.Use(m.CheckCookie)

	// TODO: move serving assets to standalone proxy like nginx
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	serveSingle("/favicon.ico", "assets/favicon.ico")
	http.Handle("/", r)

	// server
	// TODO: move port to conf
	log.Warn("starting http server on port 8000")
	err = http.ListenAndServe("0.0.0.0:8000", nil)
	if err != nil {
		log.Fatalf("starting HTTP server failed: %v", err)
	}
}

func cookieInit(hashKey, blockKey []byte) *securecookie.SecureCookie {
	return securecookie.New(hashKey, blockKey)
}

func loggingInit() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}
