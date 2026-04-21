package application

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zimlewis/shortened/global"
	"github.com/zimlewis/shortened/handler"
)

type Application struct {
	handler   http.Handler
	appConfig *global.Config
}

func New(eventCh chan []byte, appConfig *global.Config) *Application {
	app := &Application{
		appConfig: appConfig,
	}
	app.handler = app.loadHandler()

	return app
}

func (a *Application) Start() error {
	server := &http.Server{
		Handler: a.handler,
		Addr:    ":3000",
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("Cannot start server: %w", err)
	}

	return nil
}

func (app *Application) loadHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	repo := app.appConfig.Repository
	messageChannel := app.appConfig.WriteMessageChannel

	router.Get("/{id}", handler.RedirectShortened(messageChannel, repo))
	router.Get("/{id}/count", handler.GetShortenedCount(repo))
	router.Post("/", handler.AddShortened(repo))

	return router
}
