package poly

import (
  "errors"
	"params"
)

type ModularArray struct {
  n uint32
  q uint32
  data []int32
}

func NewModularArray(n,q uint32) (*ModularArray, error) {
  if n == 0 || q == 0 {
    return nil, errors.New("Invalid parameter: n or q cannot be zero")
  }
  data := make([]int32,n)
  array := &ModularArray{n,q,data}
  return array, nil
}

func (ma *ModularArray) SetData(data []int32) error {
  if ma.n != uint32(len(data)) {
    return errors.New("Mismatched data length!")
  }
  for i := 0; i < int(ma.n); i++ {
    ma.data[i] = data[i]
  }
  return nil
}

func (ma *ModularArray) GetData() []int32 {
  return ma.data
}

func (lh *ModularArray) Inc(rh *ModularArray) *ModularArray {
  if lh.n != rh.n || lh.q != rh.q {
    return nil
  }
  n,q := lh.n,lh.q
  for i := 0; i < int(n); i++ {
    lh.data[i] = addMod(lh.data[i],rh.data[i],q)
  }
  return lh
}

func (lh *ModularArray) Add(rh *ModularArray) *ModularArray {
  if lh.n != rh.n || lh.q != rh.q {
    return nil
  }
  n,q := lh.n,lh.q
  ret,_ := NewModularArray(n,q)
  for i := 0; i < int(n); i++ {
    ret.data[i] = addMod(lh.data[i],rh.data[i],q)
  }
  return ret
}

func (lh *ModularArray) Dec(rh *ModularArray) *ModularArray {
  if lh.n != rh.n || lh.q != rh.q {
    return nil
  }
  n,q := lh.n,lh.q
  for i := 0; i < int(n); i++ {
    lh.data[i] = subMod(lh.data[i],rh.data[i],q)
  }
  return lh
}

func (lh *ModularArray) Sub(rh *ModularArray) *ModularArray {
  if lh.n != rh.n || lh.q != rh.q {
    return nil
  }
  n,q := lh.n,lh.q
  ret,_ := NewModularArray(n,q)
  for i := 0; i < int(n); i++ {
    ret.data[i] = subMod(lh.data[i],rh.data[i],q)
  }
  return ret
}

func (lh *ModularArray) Mul(rh *ModularArray) *ModularArray {
  if lh.n != rh.n || lh.q != rh.q {
    return nil
  }
  n,q := lh.n,lh.q
  for i := 0; i < int(n); i++ {
    lh.data[i] = mulMod(lh.data[i],rh.data[i],q)
  }
  return lh
}

func (lh *ModularArray) Times(rh *ModularArray) *ModularArray {
  if lh.n != rh.n || lh.q != rh.q {
    return nil
  }
  n,q := lh.n,lh.q
  ret,_ := NewModularArray(n,q)
  for i := 0; i < int(n); i++ {
    ret.data[i] = mulMod(lh.data[i],rh.data[i],q)
  }
  return ret
}

func (lh *ModularArray) ScalarMul(rh int32) *ModularArray {
  n,q := lh.n,lh.q
  for i := 0; i < int(n); i++ {
    lh.data[i] = mulMod(lh.data[i],rh,q)
  }
  return lh
}

func (lh *ModularArray) ScalarTimes(rh int32) *ModularArray {
  n,q := lh.n,lh.q
  ret,_ := NewModularArray(n,q)
  for i := 0; i < int(n); i++ {
    ret.data[i] = mulMod(lh.data[i],rh,q)
  }
  return ret
}

func (lh *ModularArray) Exp(e uint32) *ModularArray {
  n,q := lh.n,lh.q
  ret,_ := NewModularArray(n,q)
  for i := 0; i < int(n); i++ {
    ret.data[i] = expMod(lh.data[i],e,q)
  }
  return ret
}

func (lh *ModularArray) bound() *ModularArray {
  n,q := lh.n,lh.q
  for i := 0; i < int(n); i++ {
    lh.data[i] = bound(lh.data[i],q)
  }
  return lh
}

func (ma *ModularArray) fft(param *params.BlissBParam) (*ModularArray,error) {
	var i,j,k uint32
	n := param.N
	q := param.Q
	psi := param.Psi
	array,err := NewModularArray(n,q)
	if err != nil {
		return nil,err
	}
	array.SetData(ma.data)
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
			if (j&k)!=0 {
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
			v[k+i] = subMod(v[k],tmp,q)
			v[k] = addMod(v[k],tmp,q)
		}
		for j = 1; j < i; j++ {
			y := psi[j * l]
			for k = j; k < n; k += i2 {
				tmp := (v[k+i] * y) % int32(q)
				v[k+i] = subMod(v[k],tmp,q)
				v[k] = addMod(v[k],tmp,q)
			}
		}
		l >>= 1
	}

	return array,nil
}

func (lh *ModularArray) flip() *ModularArray {
  n,q := lh.n,lh.q
  for i,j := 1,n-1; i < int(j); i,j = i+1,j-1 {
    tmp := lh.data[i]
    lh.data[i] = lh.data[j]
    lh.data[j] = tmp
  }
  tmp := int32(q) & ((-lh.data[0])>>31)
  lh.data[0] = tmp - lh.data[0]
  return lh
}
