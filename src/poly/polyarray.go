// Package poly provides basic manipulation of polynomials
package poly

import (
	"fmt"
	"params"
)

// PolyArray is the polynomial class provided by package poly. We don't use
// Polynomial directly as the class name, because the NTT (Numerical
// Theoretical Transform) of a polynomial is also stored as an array of n
// finite field elements. The operations of NTT are almost exactly the same
// to usual polynomial, and we don't want to implement that twice. So we just
// use PolyArray as the name, and it could be used to store either a polynomial
// or an NTT.
//
// The polynomials we deal with in BLISS are elements in the ringt Z[x]/(x^n+1)
// The modulus q is neglected in usual cases, i.e. when doing additions and
// subtractions. Only when ring multiplication and inversion are
// considered the modulus q will be considered, 'cause the NTT transform must
// be carried out in finite field, so we temporarily consider the coefficients
// as in Z_q.
type PolyArray struct {
	n     uint32
	q     uint32
	data  []int32
	param *params.BlissBParam
}

// The local poly array constructor, specified the parameters directly, create
// the array of size n, and set the parameter q. The BLISS parameter is left
// empty.
func newPolyArray(n, q uint32) (*PolyArray, error) {
	if n == 0 || q == 0 {
		return nil, fmt.Errorf("Invalid parameter: n or q cannot be zero")
	}
	data := make([]int32, n)
	array := &PolyArray{n, q, data, nil}
	return array, nil
}

// The public poly array constructor, specified the BLISS parameter set.
// Invokes the local constructor with the n and q in the BLISS parameter set.
func NewPolyArray(param *params.BlissBParam) (*PolyArray, error) {
	if param == nil {
		return nil, fmt.Errorf("Param is nil")
	}
	pa, err := newPolyArray(param.N, param.Q)
	if err != nil {
		return nil, err
	}
	pa.param = param
	return pa, err
}

// The public poly array constructor specified with the BLISS version number.
// This function looks up the parameter set by the version number, and invokes
// the other constructor with the found parameter set.
func New(version int) (*PolyArray, error) {
	param := params.GetParam(version)
	if param == nil {
		return nil, fmt.Errorf("Failed to get parameter")
	}
	return NewPolyArray(param)
}

// Return the size n of the poly array.
func (pa *PolyArray) Size() uint32 {
	return pa.n
}

// Return the BLISS parameter set of the poly array.
func (pa *PolyArray) Param() *params.BlissBParam {
	return pa.param
}

// Format the content of this array into human readable string.
func (pa *PolyArray) String() string {
	return fmt.Sprintf("%d", pa.data)
}

// Copy the given data into the content array, i.e. set the coefficient.
func (pa *PolyArray) SetData(data []int32) error {
	if pa.n != uint32(len(data)) {
		return fmt.Errorf("Mismatched data length!")
	}
	for i := 0; i < int(pa.n); i++ {
		pa.data[i] = data[i]
	}
	return nil
}

// Return the reference to the content array.
func (pa *PolyArray) GetData() []int32 {
	return pa.data
}
