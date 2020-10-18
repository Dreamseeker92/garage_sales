package handlers

import (
	"garagesale/internal/platform/web"
	"garagesale/internal/product"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	DB *gorm.DB
}

// Builder for Products handler
func NewProductsHandler(db *gorm.DB) *Products {
	return &Products{DB: db}
}

// List gets all products from the service layer and encodes them for the
// client response.
func (p *Products) List(response http.ResponseWriter, request *http.Request) error {
	list, err := product.List(p.DB)
	if err != nil {
		return errors.Wrap(err, "Listing products")
	}

	return web.Respond(response, list, http.StatusOK)
}

// Fetches a product by id.
func (p *Products) Fetch(response http.ResponseWriter, request *http.Request) error {
	id := chi.URLParam(request, "id")
	product, err := product.Fetch(p.DB, chi.URLParam(request, "id"))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(err, "record with id %s has not been found", id)
		} else {
			return errors.Wrapf(err, "error: fetching a product with  id of %s: %s", id, err)
		}

		return nil
	}

	return web.Respond(response, product, http.StatusOK)
}

// Create a new product from request parameters.
func (p *Products) Create(response http.ResponseWriter, request *http.Request) error {
	var newProduct product.Product
	if err := web.Decode(request, &newProduct); err != nil {
		return err
	}

	product, err := product.Persist(p.DB, &newProduct)
	if err != nil {
		return errors.Wrap(err, "Persisting a product")
	}

	return web.Respond(response, product, http.StatusOK)
}
