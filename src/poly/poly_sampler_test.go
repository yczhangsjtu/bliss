package poly

import (
	"params"
	"sampler"
	"testing"
)

func TestUniformPoly(t *testing.T) {
	seed := make([]uint8, sampler.SHA_512_DIGEST_LENGTH)
	for i := 0; i < len(seed); i++ {
		seed[i] = uint8(i % 8)
	}
	entropy, err := sampler.NewEntropy(seed)
	if err != nil {
		t.Errorf("Failed to create entropy")
	}
	p := UniformPoly(params.BLISS_B_4, entropy)
	if p == nil {
		t.Errorf("Failed to generate uniforma polynomial")
	}
	count0 := 0
	count1 := 0
	count2 := 0
	for i := 0; i < int(p.n); i++ {
		if p.data[i] == 0 {
			count0++
		} else if p.data[i] == 1 || p.data[i] == -1 {
			count1++
		} else if p.data[i] == 2 || p.data[i] == -2 {
			count2++
		}
	}
	if count1 != int(p.param.Nz1) {
		t.Errorf("Number of +-1 does not match: expect %d, got %d",
			p.param.Nz1, count1)
	}
	if count2 != int(p.param.Nz2) {
		t.Errorf("Number of +-2 does not match: expect %d, got %d",
			p.param.Nz2, count2)
	}
	if count0 != int(p.n-p.param.Nz1-p.param.Nz2) {
		t.Errorf("Number of 0 does not match: expect %d, got %d",
			p.n-p.param.Nz1-p.param.Nz2, count0)
	}
}
