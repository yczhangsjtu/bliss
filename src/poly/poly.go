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
	var i,j,k uint32
	n := p.param.N
	q := p.param.Q
	psi := p.param.Psi
	array,err := NewModularArray(n,q)
	if err != nil {
		return nil,err
	}
	array.SetData(p.data)
	v := array.data

	// Bit-Inverse Shuffle
	j = n >> 1
	for i = 1; i < n-1; i++ {
		if i < j {
			tmp := v[i]
			v[i] = v[j]
			v[j] = tmp
		}
		k := n
		for {
			k >>= 1
			j ^= k
			if (j&k)!=0 {
				break
			}
		}
	}

	// Main loop
	l := n
	for i = 1; i < n; i <<= 1 {
		i2 := i + i
		for k = 0; k < n; k += i2 {
			tmp := v[k+i]
			v[k+i] = subMod(v[k],tmp,q)
			v[k] = addMod(v[k],tmp,q)
		}
		for j = 1; j < i; j++ {
			y := psi[j * l]
			for k = j; k < n; k += i2 {
				tmp := (v[k+i] * y) % int32(q)
				v[k+i] = subMod(v[k],tmp,q)
				v[k] = addMod(v[k],tmp,q)
			}
		}
		l >>= 1
	}

	return array,nil
}
