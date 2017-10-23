package poly

import (
	"errors"
	"params"
)

type NTT struct {
	*ModularArray
	param *params.BlissBParam
}

func newNTT(n,q uint32) (*NTT, error) {
  array,err := NewModularArray(n,q)
  ntt := NTT{array,nil}
  if err != nil {
    return nil,err
  }
  return &ntt,err
}

func NewNTT(param *params.BlissBParam) (*NTT, error) {
	if param == nil {
		return nil,errors.New("Param is nil")
	}
	ntt,err := newNTT(param.N,param.Q)
	if err != nil {
		return nil,err
	}
	ntt.param = param
	return ntt,err
}

func (ntt *NTT) Poly() (*Polynomial,error) {
	rpsi,err := NewModularArray(ntt.n,ntt.q)
	rpsi.SetData(ntt.param.RPsi)
	if err != nil {
		return nil,err
	}
	f,err := ntt.fft(ntt.param)
	if err != nil {
		return nil,err
	}
	f.Mul(rpsi)
	f.flip()
	p := Polynomial{f,ntt.param}
	return &p,nil
}

func (ntt *NTT) Invert() (*NTT,error) {
	for i := 0; i < int(ntt.n); i++ {
		if ntt.data[i] == 0 {
			return nil,errors.New("Polynomial not invertible")
		}
	}
	ret := NTT{ntt.Exp(ntt.q-2),ntt.param}
	return &ret,nil
}
