package poly

import (
	"errors"
)

// The Fast Fourier Transform, which is the core part of NTT.
// The algorithm here is the Cooley-Tukey version of the FFT.
// For an explanation of this algorithm, refer to
// Chapter 12 of "William H. Press, etc: Numerical Recipes, 3ed".
// This function copies this polynomial and apply the FFT to the copy.
// The original polynomial is left unchanged.
// TODO: implement the local version of FFT which applies the FFT directly
//       to the original polynomial, and encapsulate the local version in the
//       copy version.
func (ma *PolyArray) FFT() (*PolyArray, error) {
	var i, j, k uint32
	n := ma.param.N
	q := ma.param.Q
	psi := ma.param.Psi
	array, err := NewPolyArray(ma.param)
	if err != nil {
		return nil, err
	}
	// Copy the original polynomial into the returned value.
	array.SetData(ma.data)
	// v is used as a reference to the polynomial coefficients, for simplicity.
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
			if (j & k) != 0 {
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
			v[k+i] = subMod(v[k], tmp, q)
			v[k] = addMod(v[k], tmp, q)
		}
		for j = 1; j < i; j++ {
			y := psi[j*l]
			for k = j; k < n; k += i2 {
				tmp := (v[k+i] * y) % int32(q)
				v[k+i] = subMod(v[k], tmp, q)
				v[k] = addMod(v[k], tmp, q)
			}
		}
		l >>= 1
	}

	return array, nil
}

// Encapsulate the FFT into the NTT procedure. NTT differentiate from FFT
// by a preprocessing procedure. In the preprocessing, multiply the i'th
// element by psi^i, where psi is sqrt(omega) mod q, where omega is a n'th
// root of unity mod q, which makes psi a 2n'th root of unity.
// TODO: implement a local version.
func (p *PolyArray) NTT() (*PolyArray, error) {
	psi, err := NewPolyArray(p.param)
	if err != nil {
		return nil, err
	}
	// Copy the parameter into a poly array.
	// p.param.Psi stores the array {psi^i}_{i=0}^{n-1}
	// TODO: Save this copy procedure by doing the loop explicitly, instead of
	//       utilizing the TimesModQ method of polyarray.
	psi.SetData(p.param.Psi)
	// Do the multiplication element-wise.
	f := p.TimesModQ(psi)
	// Apply the FFT.
	// TODO: Save the intermediate polynomial f by replacing FFT with a
	//       local version.
	g, err := f.FFT()
	if err != nil {
		return nil, err
	}
	return g, nil
}

// The Inversion NTT procedure. Instead of using IFFT, this procedure is
// carried out by still applying the FFT, based on the observation that FFT
// and IFFT are basically equivalent except for a flip.
// TODO: Implement a local version.
func (ntt *PolyArray) INTT() (*PolyArray, error) {
	// Copy the parameter into a poly array.
	// p.param.RPsi stores the array {psi^-i}_{i=0}^{n-1}
	// TODO: Save this copy procedure by doing the loop explicitly, instead of
	//       utilizing the TimesModQ method of polyarray.
	rpsi, err := NewPolyArray(ntt.param)
	rpsi.SetData(ntt.param.RPsi)
	if err != nil {
		return nil, err
	}
	f, err := ntt.FFT()
	if err != nil {
		return nil, err
	}
	f.MulModQ(rpsi)
	f.flip()
	return f, nil
}

// Invert a polynomial, assuming that the polynomial is already in NTT form.
// In NTT form, the polynomial inversion is done by an element-wise inversion
// mod q. The inversion mod q is equal to taking exponentiation q-2 mod q,
// according to Little Fermat's Theorem.
func (ntt *PolyArray) InvertAsNTT() (*PolyArray, error) {
	// Check if there is 0 element. If there is, the polynomial is noninvertible.
	for i := 0; i < int(ntt.n); i++ {
		if ntt.data[i] == 0 {
			return nil, errors.New("PolyArray not invertible")
		}
	}
	// Take the exponentiation q-2.
	ret := ntt.ExpModQ(ntt.q - 2)
	return ret, nil
}

// Multiply a polynomial by another polynomial in NTT form. The result is
// a copy in polynomial form. The original polynomial remains.
func (p *PolyArray) MultiplyNTT(ntt *PolyArray) (*PolyArray, error) {
	lh, err := p.NTT()
	if err != nil {
		return nil, err
	}
	lh.MulModQ(ntt)
	return lh.INTT()
}

// The last post-processing procedure in inverse NTT.
// Negate the first element (index 0), the 1..n-1 elements are turned around.
// The procedure is applied directly to the original polynomial, and the
// reference to this polynomial is returned.
func (lh *PolyArray) flip() *PolyArray {
	n, q := lh.n, lh.q
	// The turn-around from 1 to n-1
	for i, j := 1, n-1; i < int(j); i, j = i+1, j-1 {
		tmp := lh.data[i]
		lh.data[i] = lh.data[j]
		lh.data[j] = tmp
	}
	// The negation of the first element
	tmp := int32(q) & ((-lh.data[0]) >> 31)
	lh.data[0] = tmp - lh.data[0]
	return lh
}
