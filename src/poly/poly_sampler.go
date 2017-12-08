package poly

import (
	"sampler"
)

// Uniformly sampling a polynomial, with d1 coefficients in +-1, and d2
// coefficients of +-2. d1 and d2 are specified by the BLISS parameter set.
// The algorithm of this function is taken from the C code version of BLISS
// implementation https://github.com/SRI-CSL/bliss.
func UniformPoly(version int, entropy *sampler.Entropy) *PolyArray {
	// Create new polynomial by the version number
	p, err := New(version)
	if err != nil {
		return nil
	}

	// Take the parameter n from the polynomial
	n := p.param.N
	// Take the reference of the content
	v := p.GetData()

	// Randomly filling the +-1
	i := 0
	for i < int(p.param.Nz1) {
		x := entropy.Uint16()
		j := uint32(x>>1) % n
		mask := -(1 ^ (v[j] & 1))
		i += int(mask & 1)
		v[j] += (int32((x&1)<<1) - 1) & mask
	}

	// Randomly filling the +-2
	i = 0
	for i < int(p.param.Nz2) {
		x := entropy.Uint16()
		j := uint32(x>>1) % n
		mask := -(1 ^ ((v[j] & 1) | ((v[j] & 2) >> 1)))
		i += int(mask & 1)
		v[j] += (int32((x&1)<<2) - 2) & mask
	}

	return p
}

// Sample a random polynomial by discrete Gaussian distribution.
// All the parameters are specified by the BLISS parameter set.
// The sampler is also created from the verions number.
func GaussPoly(version int, s *sampler.Sampler) *PolyArray {
	p, err := New(version)
	if err != nil {
		return nil
	}
	n := p.param.N
	v := make([]int32, n)
	// The sampling is done by loop through the array and sampling each
	// element by one-dimensional discrete Gaussian.
	for i := 0; i < int(n); i++ {
		v[i] = s.SampleGauss()
	}
	p.SetData(v)
	return p
}

// Split the polynomial sampling procedure into the sum of two Gaussian
// polynomials. The other parameters are the same, the only difference is
// at the deviation. The splitted deviations are selected such that
//        delta_alpha^2 + delta_beta^2 \approx delta^2.
// This is the alpha version of the splitted sampling.
func GaussPolyAlpha(version int, s *sampler.Sampler) *PolyArray {
	p, err := New(version)
	if err != nil {
		return nil
	}
	n := p.param.N
	v := make([]int32, n)
	for i := 0; i < int(n); i++ {
		v[i] = s.SampleGaussCtAlpha()
	}
	p.SetData(v)
	return p
}

// Split the polynomial sampling procedure into the sum of two Gaussian
// polynomials. The other parameters are the same, the only difference is
// at the deviation. The splitted deviations are selected such that
//        delta_alpha^2 + delta_beta^2 \approx delta^2.
// This is the beta version of the splitted sampling.
func GaussPolyBeta(version int, s *sampler.Sampler) *PolyArray {
	p, err := New(version)
	if err != nil {
		return nil
	}
	n := p.param.N
	v := make([]int32, n)
	for i := 0; i < int(n); i++ {
		v[i] = s.SampleGaussCtBeta()
	}
	p.SetData(v)
	return p
}
