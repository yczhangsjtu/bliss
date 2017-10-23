package sampler

import (
	"testing"
	"io/ioutil"
	"strings"
	"params"
)

func TestSampleBer(t *testing.T) {
	p := []uint8{128,128,128,128,128,128,128,128}
	testfile,err := ioutil.ReadFile("test_data/sampler_ber_test")
	if err != nil {
		t.Errorf("Failed to open file: %s",err.Error())
	}
	filecontent := strings.TrimSpace(string(testfile))
	vs := strings.Split(filecontent," ")
	seed := make([]uint8,SHA_512_DIGEST_LENGTH)
	for i := 0; i < len(seed); i++ {
		seed[i] = uint8(i % 8)
	}
	sampler,err := New(params.BLISS_B_4,seed)
	if err != nil {
		t.Errorf("Failed to create sampler: %s",err.Error())
	}
	for i := 0; i < 512; i++ {
		bit := sampler.sampleBer(p)
		if (bit && vs[i] == "0") || (!bit && vs[i] == "1") {
			t.Errorf("Error in sampleBer: expect %s", vs[i])
		}
	}
}

func TestSampleBerExp(t *testing.T) {
	testfile,err := ioutil.ReadFile("test_data/sampler_exp_test")
	if err != nil {
		t.Errorf("Failed to open file: %s",err.Error())
	}
	filecontent := strings.TrimSpace(string(testfile))
	vs := strings.Split(filecontent," ")
	seed := make([]uint8,SHA_512_DIGEST_LENGTH)
	for i := 0; i < len(seed); i++ {
		seed[i] = uint8(i % 8)
	}
	sampler,err := New(params.BLISS_B_4,seed)
	if err != nil {
		t.Errorf("Failed to create sampler: %s",err.Error())
	}
	for i := 0; i < 512; i++ {
		bit := sampler.SampleBerExp(uint32(i * 200))
		if (bit && vs[i] == "0") || (!bit && vs[i] == "1") {
			t.Errorf("Error in sampleBerExp: expect %s", vs[i])
		}
	}
}
