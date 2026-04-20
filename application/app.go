package application

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	handler http.Handler
}

func New() *Application {
	return &Application {
		handler: loadHandler(),
	}
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

func loadHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	return router
}
