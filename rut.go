/*
Package rut validates and generates 'Rol Único Tributario'
https://en.wikipedia.org/wiki/National_identification_number#Chile
Alvaro Leiva M.
https://github.com/alvarolm
*/
package rut

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	dvseparator = '-'
)

var (
	ErrMinLength     = errors.New("length less than expected")
	ErrMaxLength     = errors.New("exceeded max length")
	ErrNoDVSeparator = errors.New("no valid 'digito verificador' separator: '-'")
	ErrInvalidDVchar = errors.New("expected digit or 'K' as 'digito verificador', instead found invalid character")
	ErrExpectedDigit = errors.New("expected digit in 'cuerpo', instead found invalid character")
	ErrinvalidDV     = errors.New("invalid 'digito verificador'")
)

// Rut implements 'Rol Único Tributario' formatting and validation
type Rut string

func (r *Rut) String() string {
	return string(*r)
}

var (
	// NNNNNNN-N
	MinRutlength = 9

	// NNNNNNNN-N
	MaxRutlength = 10
)

func NewRut(nid string) *Rut {
	rut := Rut(nid)
	return &rut
}

// format checks basic formatting constraints
// 'NNNN...-(N || K)'
func (r *Rut) format() (err error) {

	// removes point decimal points if has some
	*r = Rut(strings.Replace(string(*r), ".", "", -1))

	length := len(*r)

	if length < MinRutlength {
		return ErrMinLength
	} else if length > MaxRutlength {
		return ErrMaxLength
	}

	if string(*r)[length-2] != dvseparator {
		return ErrNoDVSeparator
	}

	dv := rune(string(*r)[length-1])
	if !unicode.IsDigit(dv) {
		switch dv {
		case 'k':
			*r = Rut(string(*r)[:length-1] + "K")
		case 'K':
			// pass
		default:
			return ErrInvalidDVchar
		}
	}

	body := string(*r)[:length-2]

	if _, err = strconv.Atoi(body); err != nil {
		return ErrExpectedDigit
	}

	return
}

type AdittionalValidationInfo struct {
	ExpectedDV rune
}

// Validate performs formatting and ecc validation (digito verificador)
func (r *Rut) Validate() (additionalinfo *AdittionalValidationInfo, err error) {

	if err = r.format(); err != nil {
		return
	}

	length := len(*r)
	body := string(*r)[:length-2]
	bodylastindex := len(body) - 1

	multsequence := []int{2, 3, 4, 5, 6, 7}
	currentmultpos := 0

	nextmult := func() (mult int) {
		mult = multsequence[currentmultpos]
		if currentmultpos == 5 {
			currentmultpos = 0
		} else {
			currentmultpos++
		}
		return
	}

	// validate
	var productssum int
	for i := range body {
		d := rune(body[bodylastindex-i])
		switch d {
		case '0':
			nextmult() // pass
		case '1':
			productssum += (1 * nextmult())
		case '2':
			productssum += (2 * nextmult())
		case '3':
			productssum += (3 * nextmult())
		case '4':
			productssum += (4 * nextmult())
		case '5':
			productssum += (5 * nextmult())
		case '6':
			productssum += (6 * nextmult())
		case '7':
			productssum += (7 * nextmult())
		case '8':
			productssum += (8 * nextmult())
		case '9':
			productssum += (9 * nextmult())
		default:
			err = ErrExpectedDigit
			return
		}
	}

	additionalinfo = &AdittionalValidationInfo{}

	switch m11 := 11 - (productssum % 11); m11 {
	case 11:
		additionalinfo.ExpectedDV = '0'
	case 10:
		additionalinfo.ExpectedDV = 'K'
	default:
		additionalinfo.ExpectedDV = rune(strconv.FormatInt(int64(m11), 10)[0])
	}

	dv := rune(string(*r)[length-1])

	if additionalinfo.ExpectedDV != dv {
		err = ErrinvalidDV
		return
	}

	return
}

// DecimalFormat returns a decimal point version
// safe to call after validation
// * panics with an unexpected format
func (r *Rut) DecimalFormat() string {
	parts := strings.Split(r.String(), string(dvseparator))
	d, _ := strconv.ParseInt(parts[0], 10, 64)
	return punto(d) + string(dvseparator) + parts[1]
}

func GenerateRut(min, max int) (rut Rut) {
	rand.Seed(time.Now().UnixNano())
	rut = Rut(strconv.Itoa(rand.Intn(max-min)+min) + "-0")
	ai, _ := rut.Validate()
	rut = Rut(string(rut)[:len(rut)-1] + string(ai.ExpectedDV))
	return
}

// helpers

func punto(v int64) string {

	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(v%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return strings.Join(parts[j:], ".")
}
