package handlers

import (
	data "command-line-arguments/home/sarpi/QuietPlace/product api/data/products.go"
	"encoding/json"
	"log"
	"net/http"
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
