package poly

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

func (lh *PolyArray) ScalarInc(rh int32) *PolyArray {
	lh.data[0] = lh.data[0] + rh
	return lh
}

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

func (lh *PolyArray) ScalarMul(rh int32) *PolyArray {
	n := lh.n
	for i := 0; i < int(n); i++ {
		lh.data[i] = lh.data[i] * rh
	}
	return lh
}

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
