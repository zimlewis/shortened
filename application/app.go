package application

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zimlewis/shortened/global"
	"github.com/zimlewis/shortened/handler"
	"github.com/zimlewis/shortened/repository"
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

	repo := repository.New(app.appConfig)

	router.Get("/{id}", handler.RedirectShortened(app.appConfig.WriteMessageChannel, repo))
	router.Post("/", handler.AddShortened(repo))

	return router
}
