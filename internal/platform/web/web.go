package web

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request) error

// Represents entry point for all web applications.
type App struct {
	mux *chi.Mux
	*log.Logger
}

// Fabric for App.
func NewApp(logger *log.Logger) *App {
	return &App{
		mux:    chi.NewRouter(),
		Logger: logger,
	}
}

// Handle connects a method and URL pattern to a particular application handler.
func (app *App) Handle(method, pattern string, handler Handler) {
	handlerFunc := func(response http.ResponseWriter, request *http.Request) {
		if err := handler(response, request); err != nil {
			body := ErrorResponse{Error: err.Error()}

			if err := Respond(response, body, http.StatusInternalServerError); err != nil {
				app.Logger.Println(err)
			}
		}

	}
	app.mux.MethodFunc(method, pattern, handlerFunc)
}

func (app *App) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	app.mux.ServeHTTP(response, request)
}
