package poly

import (
	"fmt"
	"io/ioutil"
  "testing"
	"strings"
	"strconv"
	"params"
)

func TestNewPolynomial(t *testing.T) {
  var polynomial *Polynomial
  var err error
  var n,q uint32
  n = 10
  q = 10
  polynomial,err = newPolynomial(n,q)
  if polynomial == nil || err != nil || polynomial.n != n || polynomial.q != q {
    t.Errorf("Failed to create modular polynomial for n = %d, q = %d\n",n,q)
  }

  n = 0
  q = 10
  polynomial,err = newPolynomial(n,q)
  if polynomial != nil || err == nil {
    t.Errorf("Created modular polynomial for n = %d, q = %d\n",n,q)
  }

  n = 10
  q = 0
  polynomial,err = newPolynomial(n,q)
  if polynomial != nil || err == nil {
    t.Errorf("Created modular polynomial for n = %d, q = %d\n",n,q)
  }

  n = 0
  q = 0
  polynomial,err = newPolynomial(n,q)
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

func TestFFT(t *testing.T) {
	for i := 0; i <= 4; i++ {
		for k := 0; k < 2; k++ {
			testfile,err := ioutil.ReadFile(fmt.Sprintf("test_data/fft_test_%d%d",k,i))
			if err != nil {
				t.Errorf("Failed to open file: %s",err.Error())
			}
			filecontent := strings.TrimSpace(string(testfile))
			vs := strings.Split(filecontent,"\n")
			if len(vs) != 2 {
				t.Errorf("Error in data read from test_data: len(vs) = %d",len(vs))
			}
			v1 := strings.Split(strings.TrimSpace(vs[0])," ")
			v2 := strings.Split(strings.TrimSpace(vs[1])," ")
			poly,err := New(i)
			if err != nil {
				t.Errorf("Failed to create polynomial: %s",err.Error())
			}
			if int(poly.n) != len(v1) || int(poly.n) != len(v2) {
				t.Errorf("Data size invalid: n = %d, but len(v1) = %d, len(v2) = %d",
				len(v1), len(v2))
			}
			for j := 0; j < int(poly.n); j++ {
				tmp,err := strconv.Atoi(v1[j])
				if err != nil {
					t.Errorf("Invalid integer: ",v1[j])
				}
				poly.data[j] = int32(tmp)
			}
			array,err := poly.FFT()
			if err != nil {
				t.Errorf("Error in FFT(): %s",err.Error())
			}
			for j := 0; j < int(poly.n); j++ {
				tmp,err := strconv.Atoi(v2[j])
				if err != nil {
					t.Errorf("Invalid integer: ",v2[j])
				}
				if tmp != int(array.data[j]) {
					t.Errorf("Wrong result: expect %d, got %d",tmp,array.data[j])
				}
			}
		}
	}
}
