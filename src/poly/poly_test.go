package poly

import (
  "testing"
	"params"
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


func TestNew(t *testing.T) {
  var polynomial *Polynomial
  var err error
  polynomial,err = New(params.BLISS_B_0)
  if polynomial == nil || err != nil {
		t.Errorf("Failed to create modular polynomial for BLISS_B_0: %s\n",err.Error())
  }
  polynomial,err = New(params.BLISS_B_1)
  if polynomial == nil || err != nil {
		t.Errorf("Failed to create modular polynomial for BLISS_B_1: %s\n",err.Error())
  }
  polynomial,err = New(params.BLISS_B_2)
  if polynomial == nil || err != nil {
		t.Errorf("Failed to create modular polynomial for BLISS_B_2: %s\n",err.Error())
  }
  polynomial,err = New(params.BLISS_B_3)
  if polynomial == nil || err != nil {
		t.Errorf("Failed to create modular polynomial for BLISS_B_3: %s\n",err.Error())
  }
  polynomial,err = New(params.BLISS_B_4)
  if polynomial == nil || err != nil {
		t.Errorf("Failed to create modular polynomial for BLISS_B_4: %s\n",err.Error())
  }
}
