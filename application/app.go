package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

func (a *Application) Start(ctx context.Context) error {
	server := &http.Server{
		Handler: a.handler,
		Addr:    ":3000",
	}
	defer func() {
		fmt.Println("Closing server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
		defer cancel()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			fmt.Printf("Cannot close server: %s\n", err.Error())
		}
	}() 

	errChan := make(chan error, 1)
	var err error

	go func(){
		errChan <- server.ListenAndServe()
	}()

	select {
	case <- ctx.Done():
		return nil
	case err = <- errChan:
        // ListenAndServe always returns a non-nil error
        // ErrServerClosed is expected on normal shutdown, not a real error
        if err != http.ErrServerClosed {
            return fmt.Errorf("server error: %w", err)
        }
        return nil
	}
}

func (app *Application) loadHandler() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	repo := app.appConfig.Repository
	messageChannel := app.appConfig.WriteMessageChannel

	router.Get("/{id}", handler.RedirectShortened(messageChannel, repo))
	router.Get("/{id}/count", handler.GetShortenedCount(repo))
	router.Post("/", handler.AddShortened(repo))
	router.Delete("/{id}", handler.DeleteShortened(repo))
	router.Put("/{id}", handler.UpdateShortened(repo))

	return router
}
