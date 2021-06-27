package main

import (
	"net/http"
	"github.com/johnreybacal/go-book/pkg/config"
	"github.com/johnreybacal/go-book/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	path := app.Path
	path += `/static/`
	fileServer := http.FileServer(http.Dir(path))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}