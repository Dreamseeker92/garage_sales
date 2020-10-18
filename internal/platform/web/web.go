package web

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

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
func (app *App) Handle(method, pattern string, fn http.HandlerFunc) {
	app.mux.MethodFunc(method, pattern, fn)
}

func (app *App) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	app.mux.ServeHTTP(response, request)
}
