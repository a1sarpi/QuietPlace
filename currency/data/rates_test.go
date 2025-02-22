package data

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"testing"
)

func TestRates(t *testing.T) {
	tr, err := NewRates(hclog.Default())

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Rates %#v. ", tr.rates)
}
