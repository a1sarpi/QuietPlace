package handlers

import (
	"errors"
	"github.com/a1sarpi/QuietPlace/product_api/data"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Update a products details
//
// responses:
// 	   201: noContentResponse
// 404: errorResponse
// 501: errorResponse

// Delete handles DELETE requests and remove items from the database
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)

	p.l.Debug("deleting record id", "id", id)
	err := p.productDB.DeleteProduct(id)
	if errors.Is(err, data.ErrProductNotFound) {
		p.l.Error("Unable to delete record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Error("Unable to delete record", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
