package poly

import (
	"errors"
	"params"
)

type Polynomial struct {
  *ModularArray
	param *params.BlissBParam
}

func NewPolynomial(n,q uint32) (*Polynomial, error) {
  array,err := NewModularArray(n,q)
  p := Polynomial{array,nil}
  if err != nil {
    return nil,err
  }
  return &p,err
}

func New(version int) (*Polynomial, error) {
	param := params.GetParam(version)
	if param == nil {
		return nil,errors.New("Failed to get parameter")
	}
	poly,err := NewPolynomial(param.N,param.Q)
	if err != nil {
		return nil,err
	}
	poly.param = param
	return poly,err
}
