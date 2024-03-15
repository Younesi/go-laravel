package render

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
)

type Render struct {
	Renderer   string // renderer type, go or jet
	RootPath   string
	Port       int
	Secure     bool
	ServerName string // app name
	JetViews   *jet.Set
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
}

func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(c.Renderer) {
	case "go":
		return c.GoPage(w, r, view, data)
	case "jet":
		return c.JetPage(w, r, view, variables, data)
	default:
		return errors.New("renderer not supported")
	}
}

// JetPage renders a template using the Jet templating engine
func (c *Render) JetPage(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	var vars = make(jet.VarMap)
	if variables != nil {
		vars = variables.(jet.VarMap)
	}
	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}
	t, err := c.JetViews.GetTemplate(fmt.Sprintf("%s.jet", view))
	if err != nil {
		log.Println("Error getting template : ", err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println("Error executing template : ", err)
	}

	return nil
}

// GoPage renders a template using the standard Go html/template package
func (c *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", c.RootPath, view))
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
