package render

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"

	"github.com/CloudyKit/jet/v6"
)

type Render struct {
	Renderer   string // renderer type, go or jet
	RootPath   string
	Port       int
	Secure     bool
	ServerName string // app name
	JetViews   *jet.Set
	Session    *scs.SessionManager
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Secure          bool
}

func (r *Render) LoadDefaultData(td *TemplateData, req *http.Request) *TemplateData {
	td.Secure = r.Secure
	td.CSRFToken = nosurf.Token(req)

	if r.Session.Exists(req.Context(), "UserId") {
		td.IsAuthenticated = true
	}

	return td
}

func (r *Render) Page(w http.ResponseWriter, req *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(r.Renderer) {
	case "go":
		return r.GoPage(w, req, view, data)
	case "jet":
		return r.JetPage(w, req, view, variables, data)
	}

	return errors.New("renderer not supported")
}

// GoPage renders a template using the standard Go html/template package
func (r *Render) GoPage(w http.ResponseWriter, req *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", r.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	err = tmpl.Execute(w, &td)
	if err != nil {
		return err
	}

	return nil
}

// JetPage renders a template using the Jet templating engine
func (r *Render) JetPage(w http.ResponseWriter, req *http.Request, view string, variables, data interface{}) error {
	var vars = make(jet.VarMap)
	if variables != nil {
		vars = variables.(jet.VarMap)
	}
	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}
	td = r.LoadDefaultData(td, req)

	t, err := r.JetViews.GetTemplate(fmt.Sprintf("%s.jet", view))
	if err != nil {
		log.Println("Error getting template : ", err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println("Error executing template : ", err)
	}

	return nil
}
