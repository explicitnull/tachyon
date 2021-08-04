package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"html/template"
	"net/http"
	"strings"
	"tacacs-webconsole/types"

	"github.com/sirupsen/logrus"
)

const timeShort = "2006-01-02 15:04"

// error strings for responses and logs
const (
	noIDinURL     = "no id in URL"
	databaseError = "database error"
)

const emptySelect = "--"

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

	header := types.TemplateHeader{
		Name: username,
	}

	// TODO: what is this?
	if username == "furai" {
		header.Item4 = "disabled"
	}

	hdr.Execute(w, header)
	if err != nil {
		le.WithError(err).Error("template execution failed")
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
		le.WithError(err).Error("template execution failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func executeTemplate(le *logrus.Entry, w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles("templates/" + filename)
	if err != nil {
		le.WithError(err).Error("template parsing failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
	if err != nil {
		le.WithError(err).Error("template execution failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// makeHash generates SHA hashes for given passwords
func makeHash(le *logrus.Entry, cleartext string) string {
	hash := sha256.Sum256([]byte(cleartext))
	enc := base64.StdEncoding.EncodeToString(hash[:])
	return strings.Replace(enc, "=", "", -1)
}
