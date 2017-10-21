package poly

import (
  "testing"
)

func TestNewPolynomial(t *testing.T) {
  var polynomial *Polynomial
  var err error
  var n,q uint32
  n = 10
  q = 10
  polynomial,err = NewPolynomial(n,q)
  if polynomial == nil || err != nil || polynomial.n != n || polynomial.q != q {
    t.Errorf("Failed to create modular polynomial for n = %d, q = %d\n",n,q)
  }

  n = 0
  q = 10
  polynomial,err = NewPolynomial(n,q)
  if polynomial != nil || err == nil {
    t.Errorf("Created modular polynomial for n = %d, q = %d\n",n,q)
  }

  n = 10
  q = 0
  polynomial,err = NewPolynomial(n,q)
  if polynomial != nil || err == nil {
    t.Errorf("Created modular polynomial for n = %d, q = %d\n",n,q)
  }

  n = 0
  q = 0
  polynomial,err = NewPolynomial(n,q)
  if polynomial != nil || err == nil {
    t.Errorf("Created modular polynomial for n = %d, q = %d\n",n,q)
  }
}


