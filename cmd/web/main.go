package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/johnreybacal/go-book/pkg/config"
	"github.com/johnreybacal/go-book/pkg/handlers"
	"github.com/johnreybacal/go-book/pkg/render"
	"time"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"
var app config.AppConfig

func main() {
	// Environment
	app.InProduction = false

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	
	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println("Starting application on port", portNumber)
	
	serve := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}