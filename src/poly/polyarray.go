package poly

import (
	"errors"
	"params"
	"sampler"
)

type PolyArray struct {
	n     uint32
	q     uint32
	data  []int32
	param *params.BlissBParam
}

func newPolyArray(n, q uint32) (*PolyArray, error) {
	if n == 0 || q == 0 {
		return nil, errors.New("Invalid parameter: n or q cannot be zero")
	}
	data := make([]int32, n)
	array := &PolyArray{n, q, data, nil}
	return array, nil
}

func NewPolyArray(param *params.BlissBParam) (*PolyArray, error) {
	if param == nil {
		return nil, errors.New("Param is nil")
	}
	pa, err := newPolyArray(param.N, param.Q)
	if err != nil {
		return nil, err
	}
	pa.param = param
	return pa, err
}

func New(version int) (*PolyArray, error) {
	param := params.GetParam(version)
	if param == nil {
		return nil, errors.New("Failed to get parameter")
	}
	return NewPolyArray(param)
}

func (ma *PolyArray) Size() uint32 {
	return ma.n
}

func (p *PolyArray) Param() *params.BlissBParam {
	return p.param
}

func (ma *PolyArray) SetData(data []int32) error {
	if ma.n != uint32(len(data)) {
		return errors.New("Mismatched data length!")
	}
	for i := 0; i < int(ma.n); i++ {
		ma.data[i] = data[i]
	}
	return nil
}

func (ma *PolyArray) GetData() []int32 {
	return ma.data
}

func UniformPoly(version int, entropy *sampler.Entropy) *PolyArray {
	p, err := New(version)
	if err != nil {
		return nil
	}
	n := p.param.N
	v := make([]int32, n)

	i := 0
	for i < int(p.param.Nz1) {
		x := entropy.Uint16()
		j := uint32(x>>1) % n
		mask := -(1 ^ (v[j] & 1))
		i += int(mask & 1)
		v[j] += (int32((x&1)<<1) - 1) & mask
	}

	i = 0
	for i < int(p.param.Nz2) {
		x := entropy.Uint16()
		j := uint32(x>>1) % n
		mask := -(1 ^ ((v[j] & 1) | ((v[j] & 2) >> 1)))
		i += int(mask & 1)
		v[j] += (int32((x&1)<<2) - 2) & mask
	}
	p.SetData(v)
	return p
}

func GaussPoly(version int, s *sampler.Sampler) *PolyArray {
	p, err := New(version)
	if err != nil {
		return nil
	}
	n := p.param.N
	v := make([]int32, n)
	for i := 0; i < int(n); i++ {
		v[i] = s.SampleGauss()
	}
	p.SetData(v)
	return p
}
