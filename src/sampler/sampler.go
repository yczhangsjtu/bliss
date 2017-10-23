package sampler

import (
	"fmt"
	"time"
	"math/rand"
	"params"
)


type Sampler struct {
	sigma      uint32
	ell        uint32
	prec       uint32
	columns    uint32
	kSigma     uint16
	kSigmaBits uint16

	ctable   []uint8
	cdttable []uint8

	random *rand.Rand
}

func invalidSampler() *Sampler {
	return &Sampler{0,0,0,0,0,0,[]uint8{},[]uint8{},nil}
}

func NewSampler(sigma,ell,prec uint32) (*Sampler, error) {
	columns := prec/8
	ctable,err := getTable(sigma,ell,prec)
	if err != nil {
		return invalidSampler(),err
	}
	ksigma := getKSigma(sigma,prec)
	if ksigma == 0 {
		return invalidSampler(),fmt.Errorf("Failed to get kSigma")
	}
	ksigmabits := getKSigmaBits(sigma,prec)
	if ksigmabits == 0 {
		return invalidSampler(),fmt.Errorf("Failed to get kSigmaBits")
	}
	random := rand.New(rand.NewSource(time.Now().Unix()))
	return &Sampler{sigma,ell,prec,columns,ksigma,ksigmabits,ctable,[]uint8{},random},nil
}

func New(version int) (*Sampler, error) {
	param := params.GetParam(version)
	if param == nil {
		return nil,fmt.Errorf("Failed to get parameter")
	}
	return NewSampler(param.Sigma,param.Ell,param.Prec)
}

// Sample Bernoulli distribution with probability p
// p is stored as a large big-endian integer in an array
// the real probability is p/2^d, where d is the number of
// bits of p
func (sampler *Sampler) sampleBer(p []uint8) bool {
	for _,pi := range p {
		uc := uint8(sampler.random.Uint32() & 0xff)
		if uc < pi {
			return true
		}
		if uc > pi {
			return false
		}
	}
	return true
}

// Sample Bernoulli distribution with probability p = exp(-x/(2*sigma^2))
func (sampler *Sampler) SampleBerExp(x uint32) bool {
	ri := sampler.ell - 1
	mask := uint32(1) << ri
	start := ri * sampler.columns
	for mask > 0 {
		if x & mask != 0 {
			if !sampler.sampleBer(sampler.ctable[start:start+sampler.columns]) {
				return false
			}
		}
		mask >>= 1
		start -= sampler.columns
	}
	return true
}
