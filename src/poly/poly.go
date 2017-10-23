package poly

import (
	"errors"
	"params"
)

type Polynomial struct {
  *ModularArray
	param *params.BlissBParam
}

func newPolynomial(n,q uint32) (*Polynomial, error) {
  array,err := NewModularArray(n,q)
  p := Polynomial{array,nil}
  if err != nil {
    return nil,err
  }
  return &p,err
}

func NewPolynomial(param *params.BlissBParam) (*Polynomial, error) {
	if param == nil {
		return nil,errors.New("Param is nil")
	}
  p,err := newPolynomial(param.N,param.Q)
  if err != nil {
    return nil,err
  }
	p.param = param
  return p,err
}

func New(version int) (*Polynomial, error) {
	param := params.GetParam(version)
	if param == nil {
		return nil,errors.New("Failed to get parameter")
	}
	return NewPolynomial(param)
}

func (p *Polynomial) FFT() (*ModularArray,error) {
	return p.fft(p.param)
}

func (p *Polynomial) NTT() (*NTT,error) {
	psi,err := NewModularArray(p.n,p.q)
	if err != nil {
		return nil,err
	}
	psi.SetData(p.param.Psi)
	f := Polynomial{p.Times(psi),p.param}
	g,err := f.FFT()
	if err != nil {
		return nil,err
	}
	ntt := NTT{g,p.param}
	return &ntt,nil
}
