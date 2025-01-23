package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "Water",
		Price: 1.00,
		SKU:   "abs-abc-sdf",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
