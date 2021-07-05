package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/johnreybacal/go-book/internal/config"
	"github.com/johnreybacal/go-book/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(templateData *models.TemplateData, r *http.Request) *models.TemplateData {
	templateData.CSRFToken = nosurf.Token(r)
	return templateData
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, templateData *models.TemplateData) {
	var templateCache map[string]*template.Template

	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	template, isExistent := templateCache[tmpl]
	if !isExistent {
		log.Fatal("Could not get template from template cache")
	}

	buffer := new(bytes.Buffer)
	templateData = AddDefaultData(templateData, r)

	_ = template.Execute(buffer, templateData)
	_, err := buffer.WriteTo(w)

	if err != nil {
		fmt.Println("Error parsing template: ", err)
	}
}

// CreateTemplateCache create a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	path := app.Path
	path += `\templates\`

	pages, err := filepath.Glob(path + "*.page.tmpl")

	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob(path + "*.layout.tmpl")
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(path + "*.layout.tmpl")
		}
		if err != nil {
			return cache, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}

// func RenderTemplate(w http.ResponseWriter, tmpl string) {
// 	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
// 	err = parsedTemplate.Execute(w, nil)

// 	if err != nil {
// 		fmt.Println("Error parsing template: ", err)
// 		return
// 	}
// }
