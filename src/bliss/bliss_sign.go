// Package bliss presents the core API of the BLISS signature scheme.
// Specifically, the data types for Bliss Private Key, Public Key and Signature
// and the algorithms for Key Generation, Signature Generation and Verification
package bliss

import (
	"fmt"
	"golang.org/x/crypto/sha3"
	"huffman"
	"params"
	"poly"
	"sampler"
)

// The BLISS signature structure.
// A BLISS signature contains two polynomials z1 and z2 that are bounded by
// B_inf and ceil(B_inf/2^d) respectively. The signature also contains a
// challenge c, which is an index set of size kappa in [0,n).
type BlissSignature struct {
	z1 *poly.PolyArray
	z2 *poly.PolyArray
	c  []uint32
}

// Get a human readable form of a BLISS signature.
func (sig *BlissSignature) String() string {
	return fmt.Sprintf("{z1:%s,z2:%s,c:%d}",
		sig.z1.String(), sig.z2.String(), sig.c)
}

// Compute the special Hash function from a pair of polynomial and message
// digest to a challenge, i.e. an index set of size kappa in [0,n).
// The cryptographic hash (in this case SHA3-512) of (u||hash) is used as the
// random source to generate the indices.
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
				if !array[index] {
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

// The GreedySC algorithm proposed in [Accelerating Bliss]. This algorithm
// modifies (v1,v2) := (s1,s2)*c, where c is the challenge polynomial
// represented by an index set, i.e. a sparse polynomial of coefficients 0 and 1,
// and the 1s are specified by the index set.
// GreedySC algorithm modifies the multiplication by effectively computing
// (s1,s2)*c', where c' is almost identical to c, except for the signs of a
// subset of the nonzero (i.e. equals 1) coefficients.
// c' is selected to try to minimize the norm of (s1,s2)*c'.
// This algorithm is not the most optimized, but the result is sufficiently
// satisfactory.
func greedySc(indices []uint32, s1, s2 *poly.PolyArray) (v1, v2 *poly.PolyArray) {
	n := s1.Param().N
	v1, _ = poly.NewPolyArray(s1.Param())
	v2, _ = poly.NewPolyArray(s2.Param())
	s1data := s1.GetData()
	s2data := s2.GetData()
	v1data := v1.GetData()
	v2data := v2.GetData()
	for k := 0; k < len(indices); k++ {
		index := indices[k]
		sign := int32(0)
		for i := uint32(0); i < n-index; i++ {
			sign += s1data[i]*v1data[index+i] + s2data[i]*v2data[index+i]
		}
		for i := n - index; i < n; i++ {
			sign -= s1data[i]*v1data[index+i-n] + s2data[i]*v2data[index+i-n]
		}
		if sign > 0 {
			for i := uint32(0); i < n-index; i++ {
				v1data[index+i] -= s1data[i]
				v2data[index+i] -= s2data[i]
			}
			for i := n - index; i < n; i++ {
				v1data[index+i-n] += s1data[i]
				v2data[index+i-n] += s2data[i]
			}
		} else {
			for i := uint32(0); i < n-index; i++ {
				v1data[index+i] += s1data[i]
				v2data[index+i] += s2data[i]
			}
			for i := n - index; i < n; i++ {
				v1data[index+i-n] -= s1data[i]
				v2data[index+i-n] -= s2data[i]
			}
		}
	}
	return
}

// The BLISS signature generation algorithm.
func (key *BlissPrivateKey) Sign(msg []byte, entropy *sampler.Entropy) (*BlissSignature, error) {
	kappa := key.Param().Kappa
	version := key.Param().Version
	Binf := key.Param().Binf
	Bl2 := key.Param().Bl2
	M := key.Param().M
	sampler, err := sampler.New(version, entropy)
	if err != nil {
		return nil, err
	}
	hash := sha3.Sum512(msg)
restart:
	y1 := poly.GaussPoly(version, sampler)
	y2 := poly.GaussPoly(version, sampler)
	v, err := y1.MultiplyNTT(key.a)
	if err != nil {
		return nil, err
	}
	v.ScalarMul(2)
	v.ScalarMul(int32(key.Param().OneQ2))
	v.Inc(y2)
	v = v.Mod2Q()
	dv := v.DropBits().ModP()
	indices := computeC(kappa, dv, hash[:])
	v1, v2 := greedySc(indices, key.s1, key.s2)
	normV := v1.Norm2() + v2.Norm2()
	if M <= uint32(normV) {
		return nil, fmt.Errorf("|v|^2 is larger than M")
	}
	if !sampler.SampleBerExp(M - uint32(normV)) {
		goto restart
	}
	var z1, z2 *poly.PolyArray
	b := entropy.Bit()
	if b {
		z1 = y1.Sub(v1)
		z2 = y2.Sub(v2)
	} else {
		z1 = y1.Add(v1)
		z2 = y2.Add(v2)
	}
	prodZV := z1.InnerProduct(v1) + z2.InnerProduct(v2)
	if !sampler.SampleBerCosh(prodZV) {
		goto restart
	}
	y1 = v.Sub(z2).Mod2Q().DropBits()
	v = v.DropBits()
	z2 = v.Sub(y1).BoundByP()
	if z1.MaxNorm() > int32(Binf) {
		goto restart
	}
	y2 = z2.Mul2d()
	if y2.MaxNorm() > int32(Binf) {
		goto restart
	}
	if z1.Norm2()+y2.Norm2() > int32(Bl2) {
		goto restart
	}
	return &BlissSignature{z1, z2, indices}, nil
}

// The BLISS signature generation algorithm, which is supposed to be secure
// against side-channel attacks.
func (key *BlissPrivateKey) SignAgainstSideChannel(msg []byte, entropy *sampler.Entropy) (*BlissSignature, error) {
	kappa := key.Param().Kappa
	version := key.Param().Version
	Binf := key.Param().Binf
	Bl2 := key.Param().Bl2
	M := key.Param().M
	sampler, err := sampler.New(version, entropy)
	if err != nil {
		return nil, err
	}
	hash := sha3.Sum512(msg)
restart:
	y1alpha := poly.GaussPolyAlpha(version, sampler)
	y2alpha := poly.GaussPolyAlpha(version, sampler)
	y1beta := poly.GaussPolyBeta(version, sampler)
	y2beta := poly.GaussPolyBeta(version, sampler)
	valpha, err := y1alpha.MultiplyNTT(key.a)
	vbeta, err := y1beta.MultiplyNTT(key.a)
	if err != nil {
		return nil, err
	}
	valpha.ScalarMul(2)
	vbeta.ScalarMul(2)
	valpha.ScalarMul(int32(key.Param().OneQ2))
	vbeta.ScalarMul(int32(key.Param().OneQ2))
	valpha.Inc(y2alpha)
	vbeta.Inc(y2beta)
	v := valpha.Add(vbeta)
	v = v.Mod2Q()
	dv := v.DropBits().ModP()
	indices := computeC(kappa, dv, hash[:])
	v1, v2 := greedySc(indices, key.s1, key.s2)
	normV := v1.Norm2() + v2.Norm2()
	if M <= uint32(normV) {
		return nil, fmt.Errorf("|v|^2 is larger than M")
	}
	if !sampler.SampleBerExp(M - uint32(normV)) {
		goto restart
	}
	var z1, z2 *poly.PolyArray
	b := entropy.Bit()
	if b {
		z1 = y1alpha.Sub(v1)
		z2 = y2alpha.Sub(v2)
		z1 = z1.Add(y1beta)
		z2 = z2.Add(y2beta)
	} else {
		z1 = y1alpha.Add(v1)
		z2 = y2alpha.Add(v2)
		z1 = z1.Add(y1beta)
		z2 = z2.Add(y2beta)
	}
	prodZV := z1.InnerProduct(v1) + z2.InnerProduct(v2)
	if !sampler.SampleBerCosh(prodZV) {
		goto restart
	}
	y1 := v.Sub(z2).Mod2Q().DropBits()
	v = v.DropBits()
	z2 = v.Sub(y1).BoundByP()
	if z1.MaxNorm() > int32(Binf) {
		goto restart
	}
	y2 := z2.Mul2d()
	if y2.MaxNorm() > int32(Binf) {
		goto restart
	}
	if z1.Norm2()+y2.Norm2() > int32(Bl2) {
		goto restart
	}
	return &BlissSignature{z1, z2, indices}, nil
}

// The BLISS signature verification algorithm.
func (key *BlissPublicKey) Verify(msg []byte, sig *BlissSignature) (bool, error) {
	if key.a.Param().Version != sig.z1.Param().Version {
		return false, fmt.Errorf("Mismatched signature version")
	}
	z1, z2, indices := sig.z1, sig.z2, sig.c
	param := key.a.Param()
	if z1.MaxNorm() > int32(param.Binf) {
		return false, fmt.Errorf("z1 max norm too large")
	}
	tz2 := z2.Mul2d()
	if tz2.MaxNorm() > int32(param.Binf) {
		return false, fmt.Errorf("z2 max norm too large")
	}
	if z1.Norm2()+tz2.Norm2() > int32(param.Bl2) {
		return false, fmt.Errorf("t1,z2 L2 norm too large")
	}
	hash := sha3.Sum512(msg)
	v, err := z1.MultiplyNTT(key.a)
	if err != nil {
		return false, err
	}
	v.ScalarMul(2)
	v.ScalarMul(int32(key.Param().OneQ2))
	v = v.Mod2Q()
	vdata := v.GetData()
	for i := 0; i < len(indices); i++ {
		qq := param.Q * param.OneQ2
		vdata[indices[i]] = v.NumMod2Q(vdata[indices[i]] + int32(qq))
	}
	v = v.DropBits().Add(z2).ModP()
	indicesp := computeC(param.Kappa, v, hash[:])
	for i := 0; i < len(indices); i++ {
		if indices[i] != indicesp[i] {
			return false, fmt.Errorf("Indices mismatch!")
		}
	}
	return true, nil
}

// Get the BLISS parameter set from the signature.
func (sig *BlissSignature) Param() *params.BlissBParam {
	return sig.z1.Param()
}

// Serialize the BLISS signature into binary form.
// The signature is compressed by Huffman code.
// To be accurate, part of the signature, i.e. part of s1 and the entire s2 is
// compressed.
// According to the BLISS paper, it is not suggested to compress the lower bits
// of s1, which are relatively uniformly random.
// We take the idea of StrongSwan implementation of BLISS and compress the pairs
// (s1[i]/2^8, s2[i])_{i=0}^{n-1} by Huffman codes.
// Our implementation is different from that of StrongSwan in that s1[i] is not
// nonnegative, while the s1[i]'s in StrongSwan appear to be.
// To make use of the Huffman table generation tool by StrongSwan, we compress
// by huffman coding (abs(s1[i])/2^8, s2[i])_{i=0}^{n-1} instead, and save the
// s1[i]&0xff and its sign (totally 9 bits) in byte array of size 9*n/8.
// The challenge vector c is packed in byte array of size log(n)*kappa/8.
// Finally, the entire data is prefixed by a byte specifying the BLISS version.
// The signature format is
// [ Version | low bits and sign of z1 | challenge c | huffman(z1/2^d,z2) ]
func (sig *BlissSignature) Serialize() []byte {
	cpacker := huffman.NewBitPacker()
	zpacker := huffman.NewBitPacker()
	n := sig.Param().N
	nbit := sig.Param().Nbits
	version := sig.Param().Version
	nz2 := sig.Param().Nbz2
	kappa := sig.Param().Kappa
	code := sig.Param().Code
	z1data := sig.z1.GetData()
	z2data := sig.z2.GetData()
	ret := make([]byte, 1)
	ret[0] = byte(version)
	for i := 0; i < int(kappa); i++ {
		cpacker.WriteBits(uint64(sig.c[i]), nbit)
	}
	for i := 0; i < int(n); i++ {
		bits := Abs(z1data[i]) & 0xff
		if z1data[i] < 0 {
			bits |= 0x100
		}
		zpacker.WriteBits(uint64(bits), 9)
	}
	ret = append(ret, zpacker.Data()...)
	ret = append(ret, cpacker.Data()...)
	encoder := huffman.NewHuffmanEncoder(code)
	for i := 0; i < int(n); i++ {
		z1 := Abs(z1data[i]) >> 8
		z2 := z2data[i]
		index := int(z1)*(int(nz2)*2-1) + int(z2) + int(nz2) - 1
		if index < 0 {
			fmt.Printf("z1 = %d, z2 = %d, index = %d\n", z1, z2, index)
			return []byte{}
		}
		err := encoder.Update(index)
		if err != nil {
			return []byte{}
		}
	}
	ret = append(ret, encoder.Digest()...)
	return ret
}

// Deserialize a BLISS signature from binary form.
func DeserializeBlissSignature(data []byte) (*BlissSignature, error) {
	z1, err := poly.New(int(data[0]))
	if err != nil {
		return nil, fmt.Errorf("Error in generating new polyarray: %s", err.Error())
	}
	param := z1.Param()
	z2, err := poly.NewPolyArray(param)
	if err != nil {
		return nil, fmt.Errorf("Error in generating new polyarray: %s", err.Error())
	}
	n := param.N
	kappa := param.Kappa
	nbit := param.Nbits
	// nz1 := param.Nbz1
	nz2 := param.Nbz2
	code := param.Code

	z1data := z1.GetData()
	z2data := z2.GetData()
	cdata := make([]uint32, kappa)

	csize := (nbit*kappa + 7) / 8
	lowsize := 9 * n / 8
	lowsrc := data[1 : 1+lowsize]
	csrc := data[1+lowsize : 1+lowsize+csize]
	z1z2 := data[1+lowsize+csize:]

	decoder := huffman.NewHuffmanDecoder(code, z1z2)
	zunpacker := huffman.NewBitUnpacker(lowsrc, 9*n)
	for i := 0; i < int(n); i++ {
		bits, err := zunpacker.ReadBits(9)
		if err != nil {
			return nil, fmt.Errorf("Error in unpacking lower part of z1: %s", err.Error())
		}
		sign := int32(1)
		if bits&0x100 > 0 {
			sign = int32(-1)
		}
		z1low := int32(bits & 0xff)
		index, err := decoder.Next()
		if err != nil {
			return nil, fmt.Errorf("Error in decoding huffman: %s", err.Error())
		}
		if index < 0 {
			return nil, fmt.Errorf("Invalid index %d", index)
		}
		z1high := index / (int(nz2)*2 - 1)
		z2 := int32(index%(int(nz2)*2-1) - int(nz2) + 1)
		z1 := sign * (int32(z1high<<8) | z1low)
		z1data[i] = z1
		z2data[i] = z2
	}

	cunpacker := huffman.NewBitUnpacker(csrc, nbit*kappa)
	for i := 0; i < int(kappa); i++ {
		bits, err := cunpacker.ReadBits(nbit)
		if err != nil {
			return nil, fmt.Errorf("Error in unpacking c: %s", err.Error())
		}
		cdata[i] = uint32(bits)
	}

	return &BlissSignature{z1, z2, cdata[:]}, nil
}

// A util function for computing absolute value of integers.
func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
