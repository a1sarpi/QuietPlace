package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	data "github.com/a1sarpi/QuietPlace/product api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, h *http.Request) {
	lp := data.GetProducts()
	d, err := json.Marshal(lp)

	if err != nil {
		http.ErrAbortHandler(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

	rw.Write(d)
}
