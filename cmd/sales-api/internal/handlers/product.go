package handlers

import (
	"garagesale/internal/platform/web"
	"garagesale/internal/product"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
	list, err := product.List(request.Context(), p.DB)
	if err != nil {
		return errors.Wrap(err, "Listing products")
	}

	return web.Respond(response, list, http.StatusOK)
}

// Fetches a product by id.
func (p *Products) Fetch(response http.ResponseWriter, request *http.Request) error {
	id := chi.URLParam(request, "id")
	product, err := product.Fetch(request.Context(), p.DB, id)

	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case gorm.ErrInvalidData:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting product %q", id)
		}
	}

	return web.Respond(response, product, http.StatusOK)
}

// Create a new product from request parameters.
func (p *Products) Create(response http.ResponseWriter, request *http.Request) error {
	var newProduct product.Product
	if err := web.Decode(request, &newProduct); err != nil {
		return err
	}

	product, err := product.Persist(request.Context(), p.DB, &newProduct)
	if err != nil {
		return errors.Wrap(err, "Persisting a product")
	}

	return web.Respond(response, product, http.StatusOK)
}

func (p *Products) Update(response http.ResponseWriter, request *http.Request) error {
	id := chi.URLParam(request, "id")

	var updateParams product.Product
	if err := web.Decode(request, &updateParams); err != nil {
		return err
	}

	if err := product.Update(request.Context(), p.DB, id, updateParams); err != nil {
		if err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case gorm.ErrInvalidData:
				return web.NewRequestError(err, http.StatusBadRequest)
			default:
				return errors.Wrapf(err, "Fetching product with id %s", id)
			}
		}
	}

	return web.Respond(response, nil, http.StatusNoContent)
}

func (p *Products) PersistSale(response http.ResponseWriter, request *http.Request) error {
	var sale *product.Sale

	if err := web.Decode(request, &sale); err != nil {
		return errors.Wrap(err, "Decoding new sale")
	}

	err := product.PersistSale(request.Context(), p.DB, sale)
	if err != nil {
		return errors.Wrap(err, "Adding new sale")
	}

	return web.Respond(response, sale, http.StatusCreated)
}

func (p *Products) ListSales(response http.ResponseWriter, request *http.Request) error {
	sales, err := product.ListSales(request.Context(), p.DB)
	if err != nil {
		return errors.Wrap(err, "Listing sales")
	}

	return web.Respond(response, sales, http.StatusOK)
}

func (p *Products) FetchSale(response http.ResponseWriter, request *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		return errors.Wrap(err, "Parsing url")
	}
	sale, err := product.FetchSale(request.Context(), p.DB, id)

	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case gorm.ErrInvalidData:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting sale %q", id)
		}
	}

	return web.Respond(response, sale, http.StatusOK)
}
