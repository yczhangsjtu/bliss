package poly

import (
  "testing"
)

func TestNewModularArray(t *testing.T) {
  var array *ModularArray
  var err error
  var n,q uint32
  n = 10
  q = 10
  array,err = NewModularArray(n,q)
  if array == nil || err != nil || array.n != n || array.q != q {
    t.Errorf("Failed to create modular array for n = %d, q = %d\n",n,q)
  }

  n = 0
  q = 10
  array,err = NewModularArray(n,q)
  if array != nil || err == nil {
    t.Errorf("Created modular array for n = %d, q = %d\n",n,q)
  }

  n = 10
  q = 0
  array,err = NewModularArray(n,q)
  if array != nil || err == nil {
    t.Errorf("Created modular array for n = %d, q = %d\n",n,q)
  }

  n = 0
  q = 0
  array,err = NewModularArray(n,q)
  if array != nil || err == nil {
    t.Errorf("Created modular array for n = %d, q = %d\n",n,q)
  }
}


func  TestModularArrayInc(t *testing.T) {
  fdata := []int32{0,1,2,3,4,5,6,5,4,3}
  gdata := []int32{2,4,6,1,3,5,2,1,2,3}
  hdata := []int32{2,5,1,4,0,3,1,6,6,6}
  f,_ := NewModularArray(10,7)
  g,_ := NewModularArray(10,7)
  f.SetData(fdata)
  g.SetData(gdata)
  f.Inc(g)
  res := f.GetData()
  if len(hdata) != len(res) {
    t.Errorf("Error in computing f.Inc(g): incorrect length %d", len(res))
  }
  for i := 0; i < len(hdata); i++ {
    if hdata[i] != res[i] {
      t.Errorf("Error in computing f.Inc(g): %d != %d at pos %d",
        hdata[i],res[i],i)
    }
  }
}

func  TestModularArrayAdd(t *testing.T) {
  fdata := []int32{0,1,2,3,4,5,6,5,4,3}
  gdata := []int32{2,4,6,1,3,5,2,1,2,3}
  hdata := []int32{2,5,1,4,0,3,1,6,6,6}
  f,_ := NewModularArray(10,7)
  g,_ := NewModularArray(10,7)
  f.SetData(fdata)
  g.SetData(gdata)
  h := f.Add(g)
  res := h.GetData()
  if len(hdata) != len(res) {
    t.Errorf("Error in computing f.Add(g): incorrect length %d", len(res))
  }
  for i := 0; i < len(hdata); i++ {
    if hdata[i] != res[i] {
      t.Errorf("Error in computing f.Add(g): %d != %d at pos %d",
        hdata[i],res[i],i)
    }
  }
}

func  TestModularArrayDec(t *testing.T) {
  fdata := []int32{0,1,2,3,4,5,6,5,4,3}
  gdata := []int32{2,4,6,1,3,5,2,1,2,3}
  hdata := []int32{5,4,3,2,1,0,4,4,2,0}
  f,_ := NewModularArray(10,7)
  g,_ := NewModularArray(10,7)
  f.SetData(fdata)
  g.SetData(gdata)
  f.Dec(g)
  res := f.GetData()
  if len(hdata) != len(res) {
    t.Errorf("Error in computing f.Dec(g): incorrect length %d", len(res))
  }
  for i := 0; i < len(hdata); i++ {
    if hdata[i] != res[i] {
      t.Errorf("Error in computing f.Dec(g): %d != %d at pos %d",
        hdata[i],res[i],i)
    }
  }
}

func TestModularArraySub(t *testing.T) {
  fdata := []int32{0,1,2,3,4,5,6,5,4,3}
  gdata := []int32{2,4,6,1,3,5,2,1,2,3}
  hdata := []int32{5,4,3,2,1,0,4,4,2,0}
  f,_ := NewModularArray(10,7)
  g,_ := NewModularArray(10,7)
  f.SetData(fdata)
  g.SetData(gdata)
  h := f.Sub(g)
  res := h.GetData()
  if len(hdata) != len(res) {
    t.Errorf("Error in computing f.Sub(g): incorrect length %d", len(res))
  }
  for i := 0; i < len(hdata); i++ {
    if hdata[i] != res[i] {
      t.Errorf("Error in computing f.Sub(g): %d != %d at pos %d",
        hdata[i],res[i],i)
    }
  }
}

func  TestModularArrayMul(t *testing.T) {
  fdata := []int32{0,1,2,3,4,5,6,5,4,3}
  gdata := []int32{2,4,6,1,3,5,2,1,2,3}
  hdata := []int32{0,4,5,3,5,4,5,5,1,2}
  f,_ := NewModularArray(10,7)
  g,_ := NewModularArray(10,7)
  f.SetData(fdata)
  g.SetData(gdata)
  f.Mul(g)
  res := f.GetData()
  if len(hdata) != len(res) {
    t.Errorf("Error in computing f.Mul(g): incorrect length %d", len(res))
  }
  for i := 0; i < len(hdata); i++ {
    if hdata[i] != res[i] {
      t.Errorf("Error in computing f.Mul(g): %d != %d at pos %d",
        hdata[i],res[i],i)
    }
  }
}

func  TestModularArrayScalarMul(t *testing.T) {
  fdata := []int32{0,1,2,3,4,5,6,5,4,3}
  hdata := []int32{0,4,1,5,2,6,3,6,2,5}
  f,_ := NewModularArray(10,7)
  f.SetData(fdata)
  f.ScalarMul(4)
  res := f.GetData()
  if len(hdata) != len(res) {
    t.Errorf("Error in computing f.ScalarMul(g): incorrect length %d", len(res))
  }
  for i := 0; i < len(hdata); i++ {
    if hdata[i] != res[i] {
      t.Errorf("Error in computing f.ScalarMul(g): %d != %d at pos %d",
        hdata[i],res[i],i)
    }
  }
}

func TestModularArrayScalarTimes(t *testing.T) {
  fdata := []int32{0,1,2,3,4,5,6,5,4,3}
  hdata := []int32{0,4,1,5,2,6,3,6,2,5}
  f,_ := NewModularArray(10,7)
  f.SetData(fdata)
  h := f.ScalarTimes(4)
  res := h.GetData()
  if len(hdata) != len(res) {
    t.Errorf("Error in computing f.ScalarTimes(g): incorrect length %d", len(res))
  }
  for i := 0; i < len(hdata); i++ {
    if hdata[i] != res[i] {
      t.Errorf("Error in computing f.ScalarTimes(g): %d != %d at pos %d",
        hdata[i],res[i],i)
    }
  }
}

func TestModularArrayTimes(t *testing.T) {
  fdata := []int32{0,1,2,3,4,5,6,5,4,3}
  gdata := []int32{2,4,6,1,3,5,2,1,2,3}
  hdata := []int32{0,4,5,3,5,4,5,5,1,2}
  f,_ := NewModularArray(10,7)
  g,_ := NewModularArray(10,7)
  f.SetData(fdata)
  g.SetData(gdata)
  h := f.Times(g)
  res := h.GetData()
  if len(hdata) != len(res) {
    t.Errorf("Error in computing f.Times(g): incorrect length %d", len(res))
  }
  for i := 0; i < len(hdata); i++ {
    if hdata[i] != res[i] {
      t.Errorf("Error in computing f.Times(g): %d != %d at pos %d",
        hdata[i],res[i],i)
    }
  }
}

func TestModularArrayExp(t *testing.T) {
  fdata := []int32{0,1,2,3,4,5,6,5,4,3}
  hdata := []int32{0,1,4,5,2,3,6,3,2,5}
  f,_ := NewModularArray(10,7)
  f.SetData(fdata)
  h := f.Exp(5)
  res := h.GetData()
  if len(hdata) != len(res) {
    t.Errorf("Error in computing f.Exp(5): incorrect length %d", len(res))
  }
  for i := 0; i < len(hdata); i++ {
    if hdata[i] != res[i] {
      t.Errorf("Error in computing f.Exp(5): %d != %d at pos %d",
        hdata[i],res[i],i)
    }
  }
}

func TestModularArrayFlip(t *testing.T) {
  fdata := []int32{0,1,2,3,4,5,6,5,4,3}
  hdata := []int32{0,3,4,5,6,5,4,3,2,1}
  f,_ := NewModularArray(10,7)
  f.SetData(fdata)
  f.flip()
  res := f.GetData()
  if len(hdata) != len(res) {
    t.Errorf("Error in computing f.flip(): incorrect length %d", len(res))
  }
  for i := 0; i < len(hdata); i++ {
    if hdata[i] != res[i] {
      t.Errorf("Error in computing f.flip(): %d != %d at pos %d",
        hdata[i],res[i],i)
    }
  }
}

func TestModularArrayBound(t *testing.T) {
  fdata := []int32{0,-1,-2,3,-4,5,-6,5,-4,3}
  hdata := []int32{0, 6, 5,3, 3,5, 1,5, 3,3}
  f,_ := NewModularArray(10,7)
  f.SetData(fdata)
  f.bound()
  res := f.GetData()
  if len(hdata) != len(res) {
    t.Errorf("Error in computing f.flip(): incorrect length %d", len(res))
  }
  for i := 0; i < len(hdata); i++ {
    if hdata[i] != res[i] {
      t.Errorf("Error in computing f.flip(): %d != %d at pos %d",
        hdata[i],res[i],i)
    }
  }
}
