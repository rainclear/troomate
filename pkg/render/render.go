package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/rainclear/accroo/pkg/config"
	"github.com/rainclear/accroo/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	var err error
	tc, err = CreateTemplateCache()
	if err != nil {
		log.Fatal("Could not create temmplate cache", err)
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from temmplate cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)
	n, err := buf.WriteTo(w)

	if err != nil {
		fmt.Println("Error rendering template to browser", err)
	} else {
		fmt.Println("Num of bytes written to browser: ", n)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {

	tmplCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")

	if err != nil {
		return tmplCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tmplCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return tmplCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return tmplCache, err
			}

		}

		tmplCache[name] = ts
	}

	return tmplCache, nil
}
