package handler

import (
	"html/template"
	"io"
	"net/http"
	"tachyon-web/repository"

	log "github.com/sirupsen/logrus"
)

func (g *Gateway) AppInfo(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "this is tachyon-web\n")
}

func (g *Gateway) Login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.htm")
	if err != nil {
		log.Errorf("template parsing failed: %v", err)
	}
	t.Execute(w, nil)
}

func (g *Gateway) LoginDo(w http.ResponseWriter, r *http.Request) {
		usr := r.PostFormValue("username")
		p := r.PostFormValue("password")

		dbhash, err := repository.GetPasswordHash(g.db, usr)
		if err != nil {
			fmt.Fprintf(w, "database connection failed")
		}
		if dbhash != "" {
			phash := strings.Split(dbhash, "$")
			salt := phash[2]

			cmd := exec.Command("openssl", "passwd", "-1", "-salt", salt, p)
			stdout, err := cmd.StdoutPipe()
			checkErr(err)
			cmd.Start()
			pipe := bufio.NewReader(stdout)
			line, _, err := pipe.ReadLine()
			checkErr(err)
			fhash := string(line)

			if fhash == dbhash {
				value := map[string]string{
					"name": usr,
				}
				expiration := time.Now().Add(1 * time.Hour)
				if encoded, err := sc.Encode("username", value); err == nil {
					cookie := &http.Cookie{
						Name:    "username",
						Value:   encoded,
						Path:    "/",
						Expires: expiration,
					}

					http.SetCookie(w, cookie)
				}
				mid, err := template.ParseFiles("templates/loginok.htm")
				checkErr(err)
				mid.Execute(w, nil)

				/* Logging to database */
				d := fmt.Sprintf("LOGIN-6-SUCC User %s logged in from %s", usr, r.RemoteAddr)
				errLog := logger(d, usr, usr, r.RemoteAddr)
				if errLog != nil {
					log.Printf("DB-4-LOGERROR Error writing log to db, message was: %s, error is: %s\n", d, errLog)
				}
			} else {
				t, _ := template.ParseFiles("templates/loginerror.htm")
				t.Execute(w, nil)

				d := fmt.Sprintf("LOGIN-5-FAIL Login attempt failed: wrong password for user %s from %s", usr, r.RemoteAddr)
				errLog := logger(d, usr, usr, r.RemoteAddr)
				if errLog != nil {
					log.Printf("DB-4-LOGERROR Error writing log to db, message was: %s, error is: %s", d, errLog)
				}
			}
		} else {
			//fmt.Fprintf(w, "Login not found or password incorrect")
			t, _ := template.ParseFiles("templates/loginerror.htm")
			t.Execute(w, nil)

			d := fmt.Sprintf("LOGIN-5-UKWN Login attempt failed: user %s not found, addr is %s", usr, r.RemoteAddr)
			errLog := logger(d, usr, usr, r.RemoteAddr)
			if errLog != nil {
				log.Printf("DB-4-LOGERROR Error writing log to db, message was: %s, error is: %s", d, errLog)
			}
		}
	}
}

func (g *Gateway) Logout(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "logged out!\n")
}
