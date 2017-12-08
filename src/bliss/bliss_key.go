// Package bliss presents the core API of the BLISS signature scheme.
// Specifically, the data types for Bliss Private Key, Public Key and Signature
// and the algorithms for Key Generation, Signature Generation and Verification
package bliss

import (
	"fmt"
	"huffman"
	"params"
	"poly"
	"sampler"
)

// The data structure for bliss private key.
// The two polynomials s1 and s2 are sparse polynomials with small coefficients
// the polynomial a is actually the public key stored for efficiency reason.
// a = -s2/s1
type BlissPrivateKey struct {
	s1 *poly.PolyArray
	s2 *poly.PolyArray
	a  *poly.PolyArray
}

// The data structure for bliss public key.
// a is a polynomial obtained by dividing s2 by s1 in ring Z_q[x]/(x^n+1).
// a is stored in NTT form, saving an NTT operation in the signature generation
type BlissPublicKey struct {
	a *poly.PolyArray
}

// The BLISS Key Generation procedure.
// 1. Sample polynomials f,g in Z[x]/(x^n+1). f and g are indepdent and uniform
//    in the set of polynomials with d1 coefficients of +-1 and d2 coefficients
//    equals +-2 and other coefficients 0.
// 2. Repeat step 1 until f is invertible.
// 3. The private key is (s1,s2) = (f,2g-1)
// 4. The public key is a = -s2/s1.
// The key pair is stored in a Private Key structure.
func GeneratePrivateKey(version int, entropy *sampler.Entropy) (*BlissPrivateKey, error) {
	// Generate g first, since g is not required to be invertible
	// so there is no risk of repeat the sampling.
	s2 := poly.UniformPoly(version, entropy)
	if s2 == nil {
		return nil, fmt.Errorf("Failed to generate uniform polynomial g")
	}
	// s2 = 2g-1
	s2.ScalarMul(2)
	s2.GetData()[0] -= 1

	// Prepare s2 in NTT form
	t, err := s2.NTT()
	if err != nil {
		return nil, err
	}

	for j := 0; j < 4; j++ {
		// Now we sample f
		s1 := poly.UniformPoly(version, entropy)
		if s1 == nil {
			return nil, fmt.Errorf("Failed to generate uniform polynomial f")
		}
		// Apply the NTT
		u, err := s1.NTT()
		if err != nil {
			return nil, err
		}
		u, err = u.InvertAsNTT()
		// If f is invertible, repeat the sampling.
		if err != nil {
			continue
		}
		// t = NTT(s2/s1)
		t.MulModQ(u)
		// t = INTT(NTT(s2/s1)) = s2/s1
		t, err = t.INTT()
		if err != nil {
			return nil, err
		}
		// Negate t: t = -s2/s1
		t.ScalarMulModQ(-1)
		// a = NTT(-s2/s1)
		a, err := t.NTT()
		if err != nil {
			return nil, err
		}
		key := BlissPrivateKey{s1, s2, a}
		return &key, nil
	}
	return nil, fmt.Errorf("Failed to generate invertible polynomial")
}

// Retrieve a copy of the BLISS public key from the private key.
func (privateKey *BlissPrivateKey) PublicKey() *BlissPublicKey {
	return &BlissPublicKey{privateKey.a}
}

// Retrieve the BLISS parameter set from the BLISS private key.
func (privateKey *BlissPrivateKey) Param() *params.BlissBParam {
	return privateKey.s1.Param()
}

// Retrieve the BLISS parameter set from the BLISS public key.
func (publicKey *BlissPublicKey) Param() *params.BlissBParam {
	return publicKey.a.Param()
}

// Get the human readable string of a BLISS private key.
func (privateKey *BlissPrivateKey) String() string {
	return fmt.Sprintf("{s1:%s,s2:%s,a:%s}",
		privateKey.s1.String(), privateKey.s2.String(), privateKey.a.String())
}

// Get the human readable string of a BLISS public key.
func (publicKey *BlissPublicKey) String() string {
	return fmt.Sprintf("{a:%s}", publicKey.a.String())
}

// Serialize the BLISS private key into binary form.
// The coefficients of f,g are from set {-2,-1,0,1,2}, which can be stored in
// 3 bits. So we store f=s1 and g = (s2+1)/2 in 6*n bits, and compress them
// into bytes array by a bit packer. The data is then prefixed by a byte
// specifying the BLISS version.
func (privateKey *BlissPrivateKey) Serialize() []byte {
	packer := huffman.NewBitPacker()
	n := privateKey.Param().N
	s1data := privateKey.s1.GetData()
	s2data := privateKey.s2.GetData()
	for i := 0; i < int(n); i++ {
		packer.WriteBits(uint64(s1data[i]+2), 3)
	}
	packer.WriteBits(uint64((s2data[0]+1)/2)+2, 3)
	for i := 1; i < int(n); i++ {
		packer.WriteBits(uint64(s2data[i]/2+2), 3)
	}
	ret := []byte{byte(privateKey.Param().Version)}
	return append(ret, packer.Data()...)
}

// Deserialize a BLISS private key from binary form.
func DeserializeBlissPrivateKey(data []byte) (*BlissPrivateKey, error) {
	s1, err := poly.New(int(data[0]))
	if err != nil {
		return nil, fmt.Errorf("Error in generating new polyarray: %s", err.Error())
	}
	s2, err := poly.NewPolyArray(s1.Param())
	if err != nil {
		return nil, fmt.Errorf("Error in generating new polyarray: %s", err.Error())
	}

	n := s1.Param().N
	unpacker := huffman.NewBitUnpacker(data[1:], 6*n)
	s1data := s1.GetData()
	s2data := s2.GetData()
	for i := 0; i < int(n); i++ {
		bits, err := unpacker.ReadBits(3)
		if err != nil {
			return nil, err
		}
		s1data[i] = int32(bits) - 2
	}
	bits, err := unpacker.ReadBits(3)
	if err != nil {
		return nil, err
	}
	s2data[0] = (int32(bits)-2)*2 - 1
	for i := 1; i < int(n); i++ {
		bits, err := unpacker.ReadBits(3)
		if err != nil {
			return nil, err
		}
		s2data[i] = (int32(bits) - 2) * 2
	}
	t, err := s2.NTT()
	if err != nil {
		return nil, err
	}
	u, err := s1.NTT()
	if err != nil {
		return nil, err
	}
	u, err = u.InvertAsNTT()
	if err != nil {
		return nil, err
	}
	t.MulModQ(u)
	t, err = t.INTT()
	if err != nil {
		return nil, err
	}
	t.ScalarMulModQ(-1)
	a, err := t.NTT()
	if err != nil {
		return nil, err
	}
	key := BlissPrivateKey{s1, s2, a}
	return &key, nil
}

// Serialize the Bliss Public Key into binary form.
// This is much simpler than serializing the private key, just compressing the
// bits into a byte array. The coefficients of a are uniformly random in
// [0,q), so compression is not considered. Each coefficient takes approximate
// ceil(log(q)) bits. The data is prefixed by a byte of version number.
func (publicKey *BlissPublicKey) Serialize() []byte {
	qbit := publicKey.Param().Qbits
	n := publicKey.Param().N
	packer := huffman.NewBitPacker()
	adata := publicKey.a.GetData()
	for i := 0; i < int(n); i++ {
		packer.WriteBits(uint64(adata[i]), qbit)
	}
	ret := []byte{byte(publicKey.Param().Version)}
	return append(ret, packer.Data()...)
}

// Deserialize a BLISS public key from binary form.
func DeserializeBlissPublicKey(data []byte) (*BlissPublicKey, error) {
	a, err := poly.New(int(data[0]))
	if err != nil {
		return nil, fmt.Errorf("Error in generating new polyarray: %s", err.Error())
	}
	n := a.Param().N
	qbit := a.Param().Qbits
	unpacker := huffman.NewBitUnpacker(data[1:], n*qbit)
	adata := a.GetData()
	for i := 0; i < int(n); i++ {
		bits, err := unpacker.ReadBits(qbit)
		if err != nil {
			return nil, err
		}
		adata[i] = int32(bits)
	}
	return &BlissPublicKey{a}, nil
}
