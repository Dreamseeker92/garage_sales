package handlers

import (
	"garagesale/internal/platform/web"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// API constructs a handler which knows about all API routes.
func API(db *gorm.DB, logger *log.Logger) http.Handler  {
	productsHandler := NewProductsHandler(db)
	
	app := web.NewApp(logger)
	app.Handle(http.MethodGet, "/v1/products", productsHandler.List)
	app.Handle(http.MethodGet, "/v1/products/{id}", productsHandler.Fetch)
	app.Handle(http.MethodPost, "/v1/products", productsHandler.Create)
	
	return app
}
