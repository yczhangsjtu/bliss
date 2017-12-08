// Package sampler implements all the needed cryptographic samplers needed in BLISS.
// The randomness is provided by applying SHA3-512 to an incrementing 512-bit seed.
// The sampler is deterministic with respect to the original seed.
package sampler

// Combines two consecutive bytes in a byte array into an uint16.
func combineUint16(buf []byte, i int) uint16 {
	return uint16(buf[i]) + (uint16(buf[i+1]) << 8)
}

// Combines eight consecutive bytes in a byte array into an uint64.
func combineUint64(buf []byte, i int) uint64 {
	return uint64(buf[i]) + (uint64(buf[i+1]) << 8) +
		(uint64(buf[i+2]) << 16) + (uint64(buf[i+3]) << 24) +
		(uint64(buf[i+4]) << 32) + (uint64(buf[i+5]) << 40) +
		(uint64(buf[i+6]) << 48) + (uint64(buf[i+7]) << 56)
}
