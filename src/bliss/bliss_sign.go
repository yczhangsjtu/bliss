package bliss

import (
	"golang.org/x/crypto/sha3"
	"poly"
	"sampler"
)

type BlissSignature struct {
	z1 *poly.PolyArray
	z2 *poly.PolyArray
}

func computeC(kappa uint32, u *poly.PolyArray, hash []byte) []uint32 {
	indices := make([]uint32, kappa)
	data := u.GetData()
	n := len(data)
	for i := 0; i < n; i++ {
		hash = append(hash, byte(data[i]&0xff))
		hash = append(hash, byte((data[i]>>8)&0xff))
	}
	for tries := 0; tries < 256; tries++ {
		hash[len(hash)-1]++
		whash := sha3.Sum512(hash)
		array := make([]bool, n)
		if n == 256 {
			// BLISS_B_0: we need kappa indices of 8 bits
			i := 0
			for j := 0; j < int(sampler.SHA_512_DIGEST_LENGTH); j++ {
				index := whash[j]
				if !array[j] {
					indices[i] = uint32(index)
					array[index] = true
					i++
					if i >= int(kappa) {
						return indices
					}
				}
			}
		} else {
			// BLIS_B_1234: we need kappa indices of 9 bits
			// n should be 512 now
			extraBits := byte(0)
			i := 0
			j := 0
			for j < int(sampler.SHA_512_DIGEST_LENGTH) {
				if j&7 == 0 {
					extraBits = whash[j]
					j++
				}
				index := (uint32(whash[j]) << 1) | uint32(extraBits&1)
				extraBits >>= 1
				j++
				if !array[index] {
					indices[i] = index
					array[index] = true
					i++
					if i >= int(kappa) {
						return indices
					}
				}
			}
		}
	}
	return []uint32{}
}

// func (key *BlissPrivateKey) Sign(msg []byte, entropy *sampler.Entropy) (*BlissSignature, error) {
// 	kappa := key.param.Kappa
// 	version := key.param.Version
// 	sampler, err := sampler.New(version, entropy)
// 	if err != nil {
// 		return nil, err
// 	}
// 	hash := sha3.Sum512(msg)
// 	y1 := poly.GaussPoly(version, sampler)
// 	y2 := poly.GaussPoly(version, sampler)
// 	v, err := y1.MultiplyNTT(key.a)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return nil, nil
// }
