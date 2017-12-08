// Package poly provides basic manipulation of polynomials
package poly

// Increase this polynomial by another polynomial.
// i.e. add the coefficients of another polynomial to this one element-wise.
// Return the reference to this polynomial.
func (lh *PolyArray) Inc(rh *PolyArray) *PolyArray {
	if lh.n != rh.n || lh.q != rh.q {
		return nil
	}
	n := lh.n
	for i := 0; i < int(n); i++ {
		lh.data[i] = lh.data[i] + rh.data[i]
	}
	return lh
}

// Add two polynomials and return the sum. The original polynomial remains.
func (lh *PolyArray) Add(rh *PolyArray) *PolyArray {
	if lh.n != rh.n || lh.q != rh.q {
		return nil
	}
	n, q := lh.n, lh.q
	var ret *PolyArray
	if lh.param != nil {
		ret, _ = NewPolyArray(lh.param)
	} else {
		ret, _ = newPolyArray(n, q)
	}
	for i := 0; i < int(n); i++ {
		ret.data[i] = lh.data[i] + rh.data[i]
	}
	return ret
}

// Decrease this polynomial by another polynomial.
// i.e. substract the coefficients of another polynomial from this one
// element-wise. Return the reference to this polynomial.
func (lh *PolyArray) Dec(rh *PolyArray) *PolyArray {
	if lh.n != rh.n || lh.q != rh.q {
		return nil
	}
	n := lh.n
	for i := 0; i < int(n); i++ {
		lh.data[i] = lh.data[i] - rh.data[i]
	}
	return lh
}

// Return the difference between two polynomials. The original polynomial
// remains.
func (lh *PolyArray) Sub(rh *PolyArray) *PolyArray {
	if lh.n != rh.n || lh.q != rh.q {
		return nil
	}
	n, q := lh.n, lh.q
	var ret *PolyArray
	if lh.param != nil {
		ret, _ = NewPolyArray(lh.param)
	} else {
		ret, _ = newPolyArray(n, q)
	}
	for i := 0; i < int(n); i++ {
		ret.data[i] = lh.data[i] - rh.data[i]
	}
	return ret
}

// Scale the polynomial with another, i.e. multiply each element by the
// corresponding element of another polynomial. Return the reference to the
// original polynomial.
func (lh *PolyArray) Mul(rh *PolyArray) *PolyArray {
	if lh.n != rh.n || lh.q != rh.q {
		return nil
	}
	n := lh.n
	for i := 0; i < int(n); i++ {
		lh.data[i] = lh.data[i] * rh.data[i]
	}
	return lh
}

// Returns the element-wise product of two polynomials. The original one
// remains.
func (lh *PolyArray) Times(rh *PolyArray) *PolyArray {
	if lh.n != rh.n || lh.q != rh.q {
		return nil
	}
	n, q := lh.n, lh.q
	var ret *PolyArray
	if lh.param != nil {
		ret, _ = NewPolyArray(lh.param)
	} else {
		ret, _ = newPolyArray(n, q)
	}
	for i := 0; i < int(n); i++ {
		ret.data[i] = lh.data[i] * rh.data[i]
	}
	return ret
}

// Multiply each element of the polynomial by a constant, return a reference
// to it.
func (lh *PolyArray) ScalarMul(rh int32) *PolyArray {
	n := lh.n
	for i := 0; i < int(n); i++ {
		lh.data[i] = lh.data[i] * rh
	}
	return lh
}

// Return a copy of the original polynomial with each element scaled by a
// constant.
func (lh *PolyArray) ScalarTimes(rh int32) *PolyArray {
	n, q := lh.n, lh.q
	var ret *PolyArray
	if lh.param != nil {
		ret, _ = NewPolyArray(lh.param)
	} else {
		ret, _ = newPolyArray(n, q)
	}
	for i := 0; i < int(n); i++ {
		ret.data[i] = lh.data[i] * rh
	}
	return ret
}

// Compute the L2 norm of a polynomial.
func (pa *PolyArray) Norm2() int32 {
	sum := int32(0)
	for i := 0; i < len(pa.data); i++ {
		sum += pa.data[i] * pa.data[i]
	}
	return sum
}

// Compute the Max norm of a polynomial.
func (pa *PolyArray) MaxNorm() int32 {
	max := int32(0)
	for i := 0; i < len(pa.data); i++ {
		if pa.data[i] > max {
			max = pa.data[i]
		} else if -pa.data[i] > max {
			max = -pa.data[i]
		}
	}
	return max
}

// Compute the inner product of two polynomials.
func (lh *PolyArray) InnerProduct(rh *PolyArray) int32 {
	if lh.n != rh.n || lh.q != rh.q {
		return 0
	}
	n := lh.n
	sum := int32(0)
	for i := 0; i < int(n); i++ {
		sum += lh.data[i] * rh.data[i]
	}
	return sum
}

// Return a copy of the polynomial with the lower d bits of each element cut.
// An increase of 1 is added to the higher part depending on if the most
// significant bit of the lower part is 1 or 0.
// This procedure is used in the compression of BLISS signature.
// d is specified in the BLISS parameter set.
func (pa *PolyArray) DropBits() *PolyArray {
	var ret *PolyArray
	if pa.param != nil {
		ret, _ = NewPolyArray(pa.param)
	} else {
		return nil
	}
	delta := int32(1) << pa.param.D
	halfDelta := delta >> 1
	for i := 0; i < len(pa.data); i++ {
		ret.data[i] = (pa.data[i] + halfDelta) / delta
	}
	return ret
}

// Return a copy of the original polynomial with each element multiplied by
// 2^d. d is specified in the BLISS parameter set.
// This is approximately an inversion of the DropBits procedure, and is used
// in the verification of the BLISS signature.
func (pa *PolyArray) Mul2d() *PolyArray {
	var ret *PolyArray
	if pa.param != nil {
		ret, _ = NewPolyArray(pa.param)
	} else {
		return nil
	}
	D := pa.param.D
	for i := 0; i < len(pa.data); i++ {
		ret.data[i] = pa.data[i] << D
	}
	return ret
}

// Return a copy of the original polynomial, with each element equivalent
// modulo p, and bound in [-p/2, p/2).
func (pa *PolyArray) BoundByP() *PolyArray {
	var ret *PolyArray
	if pa.param != nil {
		ret, _ = NewPolyArray(pa.param)
	} else {
		return nil
	}
	p := int32(pa.param.Modp)
	for i := 0; i < len(pa.data); i++ {
		if pa.data[i] < -p/2 {
			ret.data[i] = pa.data[i] + p
		} else if pa.data[i] >= p/2 {
			ret.data[i] = pa.data[i] - p
		} else {
			ret.data[i] = pa.data[i]
		}
	}
	return ret
}
