package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"html/template"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

const timeShort = "2006-01-02 15:04"

type Header struct {
	Name   string
	Item1  string
	Item2  string
	Item3  string
	Item4  string
	Item5  string
	Item6  string
	Item7  string
	Item8  string
	Item9  string
	Item10 string
}

func getLogger(r *http.Request) *logrus.Entry {
	ctx := r.Context()
	le := logrus.WithField("requestID", ctx.Value("requestID")).WithField("username", ctx.Value("username"))

	return le
}

func getLoggerWithoutUsername(r *http.Request) *logrus.Entry {
	ctx := r.Context()
	le := logrus.WithField("requestID", ctx.Value("requestID"))

	return le
}

func executeHeaderTemplate(le *logrus.Entry, w http.ResponseWriter, username string) {
	hdr, err := template.ParseFiles("templates/hdr.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	header := Header{
		Name: username,
	}

	// TODO: what is this?
	if username == "furai" {
		header.Item10 = "disabled"
	}

	hdr.Execute(w, header)
	if err != nil {
		le.WithError(err).Error("template parsing failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func executeFooterTemplate(le *logrus.Entry, w http.ResponseWriter) {
	ftr, err := template.ParseFiles("templates/ftr.htm")
	if err != nil {
		le.WithError(err).Error("template parsing failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ftr.Execute(w, nil)
	if err != nil {
		le.WithError(err).Error("template parsing failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func checkErr(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

// makeHash generates SHA hashes for given passwords
func makeHash(le *logrus.Entry, cleartext string) string {
	hash := sha256.Sum256([]byte(cleartext))
	enc := base64.StdEncoding.EncodeToString(hash[:])
	return strings.Replace(enc, "=", "", -1)
}

// makeFormSelect selects one column from specified table and returns it as slice
func makeFormSelect(table, col, usr string) []string {
	if col == "subdiv" {
		return []string{"europe", "asia"}
	} else if col == "prm" {
		return []string{"rw", "ro"}
	}

	out := []string{"def", "def"}

	return out
}
