package bliss

import (
	"fmt"
	"poly"
	"sampler"
	"params"
)

type BlissPrivateKey struct {
	s1 *poly.Polynomial
	s2 *poly.Polynomial
	a  *poly.NTT
	param *params.BlissBParam
}

type BlissPublicKey struct {
	a *poly.NTT
	param *params.BlissBParam
}


func GeneratePrivateKey(version int, entropy *sampler.Entropy) (*BlissPrivateKey, error) {
	// Generate g
	s2 := poly.UniformPoly(version, entropy)
	if s2 == nil {
		return nil,fmt.Errorf("Failed to generate uniform polynomial g")
	}
	// s2 = 2g-1
	s2.ScalarMul(2)
	s2.ScalarInc(-1)

	t,err := s2.NTT()
	if err != nil {
		return nil,err
	}

	for j := 0; j < 4; j++ {
		s1 := poly.UniformPoly(version, entropy)
		if s1 == nil {
			return nil,fmt.Errorf("Failed to generate uniform polynomial f")
		}
		u,err := s1.NTT()
		if err != nil {
			return nil,err
		}
		u,err = u.Invert()
		if err != nil {
			continue
		}
		t.Mul(u.ModularArray)
		apoly,err := t.Poly()
		if err != nil {
			return nil,err
		}
		apoly.ScalarMul(-1)
		a,err := apoly.NTT()
		if err != nil {
			return nil,err
		}
		key := BlissPrivateKey{s1,s2,a,s1.Param()}
		s1.Bound()
		s2.Bound()
		return &key,nil
	}
	return nil,fmt.Errorf("Failed to generate invertible polynomial")
}

func (privateKey *BlissPrivateKey) PublicKey() *BlissPublicKey {
	return &BlissPublicKey{privateKey.a,privateKey.param}
}
