package handler

import (
	"html/template"
	"net/http"
	"tachyon/types"

	"github.com/sirupsen/logrus"
)

const timeShort = "2006-01-02 15:04"

// error strings for responses and logs
const (
	noIDinURL       = "no id in URL"
	badRequest      = "bad request"
	serverError     = "server error"
	accessForbidden = "access forbidden"
)

type Notice struct {
	Title   string
	Message string
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

	header := types.TemplateHeader{
		Name: username,
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

	// TODO: return err
}
