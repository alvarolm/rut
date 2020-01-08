package rut

import (
	"fmt"
	"testing"
)

func TestRut(t *testing.T) {
	rut := Rut("11111111-1")

	if edv, err := rut.Validate(); err != nil {
		t.Error(err, string(edv.ExpectedDV))
	}

	fmt.Println("original", rut)
	fmt.Println("decimal", rut.DecimalFormat())
}

func TestGenerarRut(t *testing.T) {
	for i := 0; i < 10; i++ {
		generatedRut := GenerateRut(5000000, 23000000)
		if _, err := generatedRut.Validate(); err != nil {
			t.Error(err)
		}
		fmt.Println(generatedRut)
	}
}
