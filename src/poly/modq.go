// Package poly provides basic manipulation of polynomials
package poly

// Increase the polynomial by another one, and bound the result in [0,q) by
// modulo q.
func (lh *PolyArray) IncModQ(rh *PolyArray) *PolyArray {
	if lh.n != rh.n || lh.q != rh.q {
		return nil
	}
	n, q := lh.n, lh.q
	for i := 0; i < int(n); i++ {
		lh.data[i] = addMod(lh.data[i], rh.data[i], q)
	}
	return lh
}

// Return a copy of addition of two polynomials, with each element bound in
// [0,q) modulo q.
func (lh *PolyArray) AddModQ(rh *PolyArray) *PolyArray {
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
		ret.data[i] = addMod(lh.data[i], rh.data[i], q)
	}
	return ret
}

// Decrease the polynomial by another one, and bound the result in [0,q) modulo
// q.
func (lh *PolyArray) DecModQ(rh *PolyArray) *PolyArray {
	if lh.n != rh.n || lh.q != rh.q {
		return nil
	}
	n, q := lh.n, lh.q
	for i := 0; i < int(n); i++ {
		lh.data[i] = subMod(lh.data[i], rh.data[i], q)
	}
	return lh
}

// Return a copy of the difference of two polynomials, with each element bound
// in [0,q) modulo q.
func (lh *PolyArray) SubModQ(rh *PolyArray) *PolyArray {
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
		ret.data[i] = subMod(lh.data[i], rh.data[i], q)
	}
	return ret
}

// Scale each element by the corresponding one in another polynomial.
// Bound the result in [0,q) modulo q.
func (lh *PolyArray) MulModQ(rh *PolyArray) *PolyArray {
	if lh.n != rh.n || lh.q != rh.q {
		return nil
	}
	n, q := lh.n, lh.q
	for i := 0; i < int(n); i++ {
		lh.data[i] = mulMod(lh.data[i], rh.data[i], q)
	}
	return lh
}

// Return product of two polynomials with each element bound in [0,q) modulo q.
func (lh *PolyArray) TimesModQ(rh *PolyArray) *PolyArray {
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
		ret.data[i] = mulMod(lh.data[i], rh.data[i], q)
	}
	return ret
}

// Scale the polynomial by constant, and bound the result in [0,q) modulo q.
func (lh *PolyArray) ScalarMulModQ(rh int32) *PolyArray {
	n, q := lh.n, lh.q
	for i := 0; i < int(n); i++ {
		lh.data[i] = mulMod(lh.data[i], rh, q)
	}
	return lh
}

// Return a multiple of the polynomial, with each element bound in [0,q) modulo
// q.
func (lh *PolyArray) ScalarTimesModQ(rh int32) *PolyArray {
	n, q := lh.n, lh.q
	var ret *PolyArray
	if lh.param != nil {
		ret, _ = NewPolyArray(lh.param)
	} else {
		ret, _ = newPolyArray(n, q)
	}
	for i := 0; i < int(n); i++ {
		ret.data[i] = mulMod(lh.data[i], rh, q)
	}
	return ret
}

// Return an exponentiation of the polynomial modulo q.
func (lh *PolyArray) ExpModQ(e uint32) *PolyArray {
	n, q := lh.n, lh.q
	var ret *PolyArray
	if lh.param != nil {
		ret, _ = NewPolyArray(lh.param)
	} else {
		ret, _ = newPolyArray(n, q)
	}
	for i := 0; i < int(n); i++ {
		ret.data[i] = expMod(lh.data[i], e, q)
	}
	return ret
}

// Bound the coefficients of a polynomial in [0,q) modulo q.
func (lh *PolyArray) ModQ() *PolyArray {
	n, q := lh.n, lh.q
	for i := 0; i < int(n); i++ {
		lh.data[i] = bound(lh.data[i], q)
	}
	return lh
}

// Bound an integer in [0,q) modulo q.
// This is implemented as a method of polyarray to save the effort of
// retrieving the modulo q of this polynomial from outside this package.
func (lh *PolyArray) NumModQ(a int32) int32 {
	return bound(a, lh.q)
}

// Bound an integer in [0,2q) modulo 2q.
// This is used in BLISS signature generation algorithm.
func (lh *PolyArray) NumMod2Q(a int32) int32 {
	return bound(a%int32(lh.q*2), lh.q*2)
}

// Bound the coefficients of a polynomial in [0,2q) modulo 2q.
func (lh *PolyArray) Mod2Q() *PolyArray {
	n := lh.n
	for i := 0; i < int(n); i++ {
		lh.data[i] = lh.NumMod2Q(lh.data[i])
	}
	return lh
}

// Bound the coefficients of a polynomial in [0,p) modulo p.
func (lh *PolyArray) ModP() *PolyArray {
	n := lh.n
	for i := 0; i < int(n); i++ {
		lh.data[i] = bound(lh.data[i]%int32(lh.param.Modp), lh.param.Modp)
	}
	return lh
}
