package handlers

import (
	"encoding/json"
	"errors"
	"garagesale/internal/product"
	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	DB *gorm.DB
	*log.Logger
}

// Builder for Products handler
func NewProductsHandler(db *gorm.DB, logger *log.Logger) *Products {
	return &Products{DB: db, Logger: logger}
}

// List gets all products from the service layer and encodes them for the
// client response.
func (p *Products) List(w http.ResponseWriter, r *http.Request) {
	list, err := product.List(p.DB)
	if err != nil {
		p.Logger.Printf("error: listing products: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		p.Logger.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		p.Logger.Println("error writing result", err)
	}
}

// Fetches a product by id.
func (p *Products) Fetch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product, err := product.Fetch(p.DB, chi.URLParam(r, "id"))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Logger.Printf("record with id %s has not been found", id)
			w.WriteHeader(http.StatusNotFound)
		} else {
			p.Logger.Printf("error: fetching a product with  id of %s: %s", id, err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	data, err := json.Marshal(product)
	if err != nil {
		p.Logger.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		p.Logger.Println("error writing result", err)
	}
}
