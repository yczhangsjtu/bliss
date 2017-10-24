package poly

import (
	"fmt"
	"io/ioutil"
  "testing"
	"strings"
	"strconv"
)

func TestInvert(t *testing.T) {
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
				t.Errorf("Invalid integer: %s",v2[j])
			}
			if tmp != int(ntt.data[j]) {
				t.Errorf("Wrong result of FFT(): expect %d, got %d",tmp,ntt.data[j])
			}
		}
		inv,err := ntt.Invert()
		if err == nil {
			test := inv.Times(ntt.ModularArray)
			for j := 0; j < int(test.n); j++ {
				if test.data[j] != test.data[j] {
					t.Errorf("Wrong result of Invert(): expect 1, got %d",test.data[j])
				}
			}
		} else {
			fmt.Printf("Test polynomial test_data/ntt_test_%d not invertible.\n",i)
		}
	}
}
