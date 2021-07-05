package handler

import (
	"bufio"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"tachyon-web/repository"

	"github.com/gorilla/mux"
)

type FormOptions struct {
	Opt1 []string
	Opt2 []string
}

func (g *Gateway) CreateEntity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username := ctx.Value("username").(string)

	if repository.CheckRole(le, g.db, username) != "admin" {
		le.Warn("access forbidden")
		fmt.Fprintf(w, "access forbidden")
		return
	}

	vars := mux.Vars(r)
	entity := vars["entity"]

	header := Header{
		Name: username,
	}

	hdr, err := template.ParseFiles("templates/hdr.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
		return
	}
	hdr.Execute(w, header)

	if entity == "user" {
		form := new(FormOptions)
		form.Opt1 = makeFormSelect("subdiv", "subdiv", "none")
		form.Opt2 = makeFormSelect("prm", "prm", "none")

		t, err := template.ParseFiles("templates/newuser.htm")
		if err != nil {
			le.WithError(err).Error("template parsing failed")
			return
		}

		err = t.Execute(w, form)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Template not found: %v", err)
			log.Printf("Template not found: %v", err)
		}

		ftr, err := template.ParseFiles("templates/ftr.htm")
		if err != nil {
			le.WithError(err).Error("template parsing failed")
			return
		}
		ftr.Execute(w, nil)
	} else if entity == "permission" {
		t, err := template.ParseFiles("templates/newpermission.htm")
		if err != nil {
			le.WithError(err).Error("template parsing failed")
			return
		}

		err = t.Execute(w, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Template not found: %v", err)
			log.Printf("Template not found: %v", err)
		}

		ftr, err := template.ParseFiles("templates/ftr.htm")
		if err != nil {
			le.WithError(err).Error("template parsing failed")
			return
		}
		ftr.Execute(w, nil)
	}
}

func (g *Gateway) CreateEntityDo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := getLogger(r)

	username := ctx.Value("username").(string)

	if repository.CheckRole(le, g.db, username) != "admin" {
		le.Warn("access forbidden")
		fmt.Fprintf(w, "access forbidden")
		return
	}

	vars := mux.Vars(r)
	entity := vars["entity"]

	header := Header{
		Name: username,
	}

	hdr, err := template.ParseFiles("templates/hdr.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
		return
	}
	hdr.Execute(w, header)

	r.ParseForm()

	if entity == "user" {
		u := new(user)

		u.Name = r.PostFormValue("username")

		u.Mail = r.PostFormValue("mail")
		if u.Mail == "" {
			u.Mail = "unknown"
		}

		u.Subdiv = r.PostFormValue("subdiv")
		u.Prm = r.PostFormValue("perm")

		/* Get IDs for saving in tables for normalization */
		subdiv_id := repository.GetDivID(le, g.db, u.Subdiv)
		prm_id, err := repository.GetPermId(le, g.db, u.Prm)
		if err != nil {
			le.WithError(err).Error("getting permissions failed")
			return
		}

		u.Cleartext = genPass()
		u.Hash = makeHash(le, u.Cleartext)

		repository.CreateUser
		if errEx != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Template error for user %s: %s", authUser, err)
			return
		} else {
			/* Logging to database */

			mid, err := template.ParseFiles("templates/usercreated.htm")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("HTTP error for user %s: %s", authUser, err)
				return
			}
			mid.Execute(w, u)

			ftr, _ := template.ParseFiles("templates/ftr.htm")
			ftr.Execute(w, nil)
		}
	} else if entity == "permission" {
		p := new(Permission)
		/*
			f1 := r.Form["perm"]
			p.Name = f1[0]
			f2 := r.Form["descr"]
			p.Descr = f2[0]
		*/
		p.Name = r.PostFormValue("perm")
		p.Descr = r.PostFormValue("descr")

		db, err := sql.Open("postgres", dbconf())
		checkErr(err)
		defer db.Close()

		err = db.Ping()
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			log.Printf("DB-2-CONERROR %s", err)
		} else {
			stmt, err := db.Prepare("INSERT INTO prm(prm, comment, created_by) values ($1, $2, $3)")
			checkErr(err)
			defer stmt.Close()

			_, errEx := stmt.Exec(p.Name, p.Descr, authUser)
			if errEx != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("DB-3-INSERR error for user %s: %s", authUser, err)
			} else {
				mid, err := template.ParseFiles("templates/permissioncreated.htm")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					log.Printf("HTTP error for user %s: %s", authUser, err)
					return
				}
				mid.Execute(w, p)
			}

			ftr, _ := template.ParseFiles("templates/ftr.htm")
			ftr.Execute(w, nil)
		}
	}
}

func makeFormSelect(table, col, usr string) []string {
	/* This function selects one column from specified table and returns it as slice */
	db, err := sql.Open("postgres", dbconf())
	checkErr(err)
	defer db.Close()

	err = db.Ping()
	checkErr(err)

	q := fmt.Sprintf("SELECT %s FROM %s WHERE act='true' ORDER by %s", col, table, col)
	//log.Println("DEBUG: q =", q)

	rows, err := db.Query(q)
	checkErr(err)
	defer rows.Close()

	var option string
	swap := make([]string, 500)

	i := 0
	swap[i] = "--"
	for rows.Next() {
		i++
		err = rows.Scan(&option)
		checkErr(err)
		swap[i] = option
	}
	nonNullCount := 0

	for i := 0; i < len(swap); i++ {
		if swap[i] != "" {
			nonNullCount++
		}
	}

	out := make([]string, nonNullCount)
	copy(out, swap)

	err = rows.Err()
	checkErr(err)
	return out
}

func genPass() string {
	cmd := exec.Command("openssl", "rand", "-base64", "7")
	stdout, err := cmd.StdoutPipe()
	checkErr(err)

	cmd.Start()
	pipe := bufio.NewReader(stdout)

	line, _, err := pipe.ReadLine()
	checkErr(err)
	res := string(line)
	pass := strings.Replace(res, "=", "", 2)
	return pass
}
