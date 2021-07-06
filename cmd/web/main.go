package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/johnreybacal/go-book/internal/config"
	"github.com/johnreybacal/go-book/internal/handlers"
	"github.com/johnreybacal/go-book/internal/models"
	"github.com/johnreybacal/go-book/internal/render"
)

const portNumber = ":8080"
var app config.AppConfig

func main() {
	initApp()

	fmt.Println("Starting application on port", portNumber)
	
	serve := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	err := serve.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}

func initApp() {
	app.InProduction = false
	app.UseCache = false
	render.NewTemplates(&app)
	initPath()
	initSession()
	initTemplateCache()
	initRepository()
}

func initSession() {
	gob.Register(models.Reservation{})
	
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	
	app.Session = session
}

func initPath() {
	path, _ := os.Getwd()

	// is launched from .vscode/launch.json
	if strings.Contains(path, `\cmd\web`) {
		path = path[:len(path) - 8]
	}
	app.Path = path
	fmt.Println("Root:", path)
}

func initTemplateCache() {
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = templateCache
}

func initRepository() {
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
}