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
