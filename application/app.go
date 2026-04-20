package application

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zimlewis/shortened/handler"
)

type Application struct {
	handler http.Handler
	Channel      chan []byte
}

func New(eventCh chan []byte) *Application {
	app := &Application {
		Channel: eventCh,
	}
	app.handler = app.loadHandler();

	return app
}

func (a *Application) Start() error {
	server := &http.Server {
		Handler: a.handler,
		Addr: ":3000",
	}

	err := server.ListenAndServe();
	if err != nil {
		return fmt.Errorf("Cannot start server: %w", err)
	}

	return nil
}

func (app *Application) loadHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/{id}", handler.RedirectShortened(app.Channel))

	return router
}
