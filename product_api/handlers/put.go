package handlers

import (
	"errors"
	"net/http"

	"github.com/a1sarpi/QuietPlace/product_api/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//		201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update products
func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Debug("updating record id", "id", prod.ID)

	err := p.productDB.UpdateProduct(prod)
	if errors.Is(err, data.ErrProductNotFound) {
		p.l.Error("Unable to find product", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, rw)
		return
	}

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}
