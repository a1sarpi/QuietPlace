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

	p.l.Println("[DEBUG] deleting record id", id)
	err := data.DeleteProduct(id)
	if errors.Is(err, data.ErrProductNotFound) {
		p.l.Println("[ERROR] deleting record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Println("[ERROR] deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
