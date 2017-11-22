package main

import (
	"fmt"
	"params"
	"sampler"
)

func main() {
	seed := make([]uint8, sampler.SHA_512_DIGEST_LENGTH)
	for i := 0; i < len(seed); i++ {
		seed[i] = uint8(i % 8)
	}
	entropy, err := sampler.NewEntropy(seed)
	if err != nil {
		fmt.Errorf("Failed to create entropy: %s", err.Error())
	}
	s, err := sampler.New(params.BLISS_B_4, entropy)
	if err != nil {
		fmt.Errorf("Failed to create sampler: %s", err.Error())
	}
	for i := 0; i < 10240; i++ {
		fmt.Printf("%d ", s.SampleGaussCtAlpha())
		// if i > 0 && i%16 == 0 {
		// 	fmt.Printf("\n")
		// }
	}
	fmt.Printf("\n")
}
