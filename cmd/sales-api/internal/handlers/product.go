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
func (p *Products) List(response http.ResponseWriter, request *http.Request) {
	list, err := product.List(p.DB)
	if err != nil {
		p.Logger.Printf("error: listing products: %s", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		p.Logger.Println("error marshalling result", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	if _, err := response.Write(data); err != nil {
		p.Logger.Println("error writing result", err)
	}
}

// Fetches a product by id.
func (p *Products) Fetch(response http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	product, err := product.Fetch(p.DB, chi.URLParam(request, "id"))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.Logger.Printf("record with id %s has not been found", id)
			response.WriteHeader(http.StatusNotFound)
		} else {
			p.Logger.Printf("error: fetching a product with  id of %s: %s", id, err)
			response.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	data, err := json.Marshal(product)
	if err != nil {
		p.Logger.Println("error marshalling result", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	if _, err := response.Write(data); err != nil {
		p.Logger.Println("error writing result", err)
	}
}

func (p *Products) Create(response http.ResponseWriter, request *http.Request) {
	var newProduct product.Product
	if err := json.NewDecoder(request.Body).Decode(&newProduct); err != nil {
		p.Logger.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := product.Persist(p.DB, &newProduct)
	if err != nil {
		p.Logger.Println("error persisting a product", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(product)
	if err != nil {
		p.Logger.Println("error marshalling result", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	if _, err := response.Write(data); err != nil {
		p.Logger.Println("error writing result", err)
	}
}
