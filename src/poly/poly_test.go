package poly

import (
	"fmt"
	"io/ioutil"
  "testing"
	"strings"
	"strconv"
	"params"
	"sampler"
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

func TestNTT(t *testing.T) {
	for i := 0; i <= 4; i++ {
		testfile,err := ioutil.ReadFile(fmt.Sprintf("test_data/ntt_test_%d",i))
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
		poly.Bound()
		ntt,err := poly.NTT()
		if err != nil {
			t.Errorf("Error in FFT(): %s",err.Error())
		}
		for j := 0; j < int(poly.n); j++ {
			tmp,err := strconv.Atoi(v2[j])
			if err != nil {
				t.Errorf("Invalid integer: ",v2[j])
			}
			if tmp != int(ntt.data[j]) {
				t.Errorf("Wrong result of FFT(): expect %d, got %d",tmp,ntt.data[j])
			}
		}
		npoly,err := ntt.Poly()
		if err != nil {
			t.Errorf("Error in Poly(): %s",err.Error())
		}
		for j := 0; j < int(poly.n); j++ {
			if npoly.data[j] != poly.data[j] {
				t.Errorf("Wrong result of Poly(): expect %d, got %d",poly.data[j],npoly.data[j])
			}
		}
	}
}

func TestUniformPoly(t *testing.T) {
	seed := make([]uint8,sampler.SHA_512_DIGEST_LENGTH)
	for i := 0; i < len(seed); i++ {
		seed[i] = uint8(i % 8)
	}
	entropy,err := sampler.NewEntropy(seed)
	if err != nil {
		t.Errorf("Failed to create entropy")
	}
	p := UniformPoly(params.BLISS_B_4, entropy)
	if p == nil {
		t.Errorf("Failed to generate uniforma polynomial")
	}
	count0 := 0
	count1 := 0
	count2 := 0
	for i := 0; i < int(p.n); i++ {
		if p.data[i] == 0 {
			count0++
		} else if p.data[i] == 1 || p.data[i] == -1 {
			count1++
		} else if p.data[i] == 2 || p.data[i] == -2 {
			count2++
		}
	}
	if count1 != int(p.param.Nz1) {
		t.Errorf("Number of +-1 does not match: expect %d, got %d",
			p.param.Nz1,count1)
	}
	if count2 != int(p.param.Nz2) {
		t.Errorf("Number of +-2 does not match: expect %d, got %d",
			p.param.Nz2,count2)
	}
	if count0 != int(p.n-p.param.Nz1-p.param.Nz2) {
		t.Errorf("Number of 0 does not match: expect %d, got %d",
			p.n-p.param.Nz1-p.param.Nz2,count0)
	}
}
