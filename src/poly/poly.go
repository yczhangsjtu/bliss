package poly

import (
	"errors"
	"sampler"
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

func UniformPoly(version int, entropy *sampler.Entropy) *Polynomial {
	p,err := New(version)
	if err != nil {
		return nil
	}
	n := p.param.N
	v := make([]int32,n)

	i := 0
	for i < int(p.param.Nz1) {
		x := entropy.Uint16()
		j := uint32(x >> 1) % n
		mask := -(1^(v[j]&1))
		i += int(mask&1)
		v[j] += (int32((x&1)<<1)-1)&mask
	}

	i = 0
	for i < int(p.param.Nz2) {
		x := entropy.Uint16()
		j := uint32(x >> 1) % n
		mask := -(1^((v[j]&1)|((v[j]&2)>>1)))
		i += int(mask&1)
		v[j] += (int32((x&1)<<2)-2)&mask
	}
	p.SetData(v)
	return p
}
