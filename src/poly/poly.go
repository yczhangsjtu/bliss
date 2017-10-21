package poly

type Polynomial struct {
  *ModularArray
}

func NewPolynomial(n,q uint32) (*Polynomial, error) {
  array,err := NewModularArray(n,q)
  p := Polynomial{array}
  if err != nil {
    return nil,err
  }
  return &p,err
}
