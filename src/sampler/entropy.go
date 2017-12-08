// Package sampler implements all the needed cryptographic samplers needed in BLISS.
// The randomness is provided by applying SHA3-512 to an incrementing 512-bit seed.
// The sampler is deterministic with respect to the original seed.
package sampler

import (
	"fmt"
	"golang.org/x/crypto/sha3"
)

// The constants used by the Entropy class.
const (
	// The number of bytes of a SHA3-512 digest.
	SHA_512_DIGEST_LENGTH uint32 = 64
	// The number of hash digests to calculate in each round.
	// Each pool is refreshed by this many hash digests!
	EPOOL_HASH_COUNT = 10
	// This many random chars can be provided in each round
	CHAR_POOL_SIZE = SHA_512_DIGEST_LENGTH * EPOOL_HASH_COUNT
	// This many random int16 can be provided in each round
	INT16_POOL_SIZE = SHA_512_DIGEST_LENGTH / 2 * EPOOL_HASH_COUNT
	// This many random int64 can be provided in each round
	INT64_POOL_SIZE = SHA_512_DIGEST_LENGTH / 8 * EPOOL_HASH_COUNT
)

// The Entropy class encapsulates the seed and SHA3-512 procedure, and exposes
// random strings of bits, char, int16 or int64.
type Entropy struct {
	bitpool   uint64
	charpool  []uint8
	int16pool []uint16
	int64pool []uint64
	seed      []uint8

	bitp   uint32
	charp  uint32
	int16p uint32
	int64p uint32
}

// Create a new instance of the Entropy class, providing the initial seed,
// which completely determines the sampling results of this entropy (provided
// that the sampling procedure is the same, for example, if you sample
// two bits, then three chars, and 10 int64, then you must do as exactly in
// the same manner to ensure that the sampling result is identical given the
// same seed.).
func NewEntropy(seed []uint8) (*Entropy, error) {
	if len(seed) < int(SHA_512_DIGEST_LENGTH) {
		return nil, fmt.Errorf("Insufficient seed length, need %d, got %d",
			SHA_512_DIGEST_LENGTH, len(seed))
	}
	entropy := Entropy{0, []uint8{}, []uint16{}, []uint64{}, []uint8{}, 0, 0, 0, 0}
	entropy.charpool = make([]uint8, CHAR_POOL_SIZE)
	entropy.int16pool = make([]uint16, INT16_POOL_SIZE)
	entropy.int64pool = make([]uint64, INT64_POOL_SIZE)
	entropy.seed = make([]uint8, SHA_512_DIGEST_LENGTH)
	for i := 0; i < int(SHA_512_DIGEST_LENGTH); i++ {
		entropy.seed[i] = seed[i]
	}
	entropy.refreshCharPool()
	entropy.refreshInt16Pool()
	entropy.refreshInt64Pool()
	entropy.refreshBitPool()
	return &entropy, nil
}

// Increase the seed as if it is a big number stored in little endian.
func (entropy *Entropy) incrementSeed() {
	for i := 0; i < int(SHA_512_DIGEST_LENGTH); i++ {
		entropy.seed[i]++
		if entropy.seed[i] > 0 {
			break
		}
	}
}

// Refresh the char pool when they are running out, by incrementing the seed
// and filling the pool with hash digests.
func (entropy *Entropy) refreshCharPool() {
	for i := 0; i < int(EPOOL_HASH_COUNT); i++ {
		offset := i * int(SHA_512_DIGEST_LENGTH)
		sha := sha3.Sum512([]byte(entropy.seed))
		for j := 0; j < int(SHA_512_DIGEST_LENGTH); j++ {
			entropy.charpool[offset+j] = uint8(sha[j])
		}
		entropy.incrementSeed()
	}
	entropy.charp = 0
}

// Refresh the int16 pool when they are running out, by incrementing the seed
// and filling the pool with hash digests.
func (entropy *Entropy) refreshInt16Pool() {
	for i := 0; i < int(EPOOL_HASH_COUNT); i++ {
		offset := i * int(SHA_512_DIGEST_LENGTH) / 2
		sha := sha3.Sum512([]byte(entropy.seed))
		for j := 0; j < int(SHA_512_DIGEST_LENGTH)/2; j++ {
			entropy.int16pool[offset+j] = combineUint16(sha[:], j*2)
		}
		entropy.incrementSeed()
	}
	entropy.int16p = 0
}

// Refresh the int64 pool when they are running out, by incrementing the seed
// and filling the pool with hash digests.
func (entropy *Entropy) refreshInt64Pool() {
	for i := 0; i < int(EPOOL_HASH_COUNT); i++ {
		offset := i * int(SHA_512_DIGEST_LENGTH) / 8
		sha := sha3.Sum512([]byte(entropy.seed))
		for j := 0; j < int(SHA_512_DIGEST_LENGTH)/8; j++ {
			entropy.int64pool[offset+j] = combineUint64(sha[:], j*8)
		}
		entropy.incrementSeed()
	}
	entropy.int64p = 0
}

// Refresh the bit pool when they are running out, by replacing the pool with
// the next random uint64.
func (entropy *Entropy) refreshBitPool() {
	entropy.bitpool = entropy.Uint64()
	entropy.bitp = 0
}

// Get the next random uint64.
func (entropy *Entropy) Uint64() uint64 {
	if entropy.int64p >= INT64_POOL_SIZE {
		entropy.refreshInt64Pool()
	}
	ret := entropy.int64pool[entropy.int64p]
	entropy.int64p++
	return ret
}

// Get the next random uint16.
func (entropy *Entropy) Uint16() uint16 {
	if entropy.int16p >= INT16_POOL_SIZE {
		entropy.refreshInt16Pool()
	}
	ret := entropy.int16pool[entropy.int16p]
	entropy.int16p++
	return ret
}

// Get the next random char.
func (entropy *Entropy) Char() uint8 {
	if entropy.charp >= CHAR_POOL_SIZE {
		entropy.refreshCharPool()
	}
	ret := entropy.charpool[entropy.charp]
	entropy.charp++
	return ret
}

// Get the next random bit.
func (entropy *Entropy) Bit() bool {
	if entropy.bitp >= 64 {
		entropy.refreshBitPool()
	}
	bit := entropy.bitpool & 1
	entropy.bitpool >>= 1
	entropy.bitp++
	return bit == 1
}

// Get a specific number of random bits packed in the lower n bits of
// an uint32.
func (entropy *Entropy) Bits(n int) uint32 {
	ret := uint32(0)
	for n > 0 {
		ret <<= 1
		if entropy.Bit() {
			ret |= 1
		}
		n--
	}
	return ret
}
