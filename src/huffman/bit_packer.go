// Package huffman implements the huffman encoding and decoding scheme.
package huffman

import (
	"fmt"
)

const (
	// Maximum number of BYTES that can be packed in the bit packer.
	// The number of bits is 8 times this.
	MAX_PACKER_SIZE int = 8192
)

// The BitPacker class implements a bit packer, which starts up empty, and
// stores arbitrary number of bits into a byte array.
// It keeps track of the number of bits stored, and can at any time exposes the
// slice of bytes that is storing the bits, which can input into and be
// unpacked by a bit unpacker.
type BitPacker struct {
	data  [MAX_PACKER_SIZE]byte
	nbyte uint32
	nbit  uint32
}

// The BitUnpacker class unpacks a byte array into arbitrary size consecutive
// bits. It starts up with a full byte array, and the bit pointer pointing to
// position zero.
type BitUnpacker struct {
	data  []byte
	pbyte uint32
	pbit  uint32
	size  uint32
}

// Create a new bit packer.
func NewBitPacker() *BitPacker {
	return &BitPacker{
		[MAX_PACKER_SIZE]byte{}, 0, 0,
	}
}

// Create a new bit unpacker to unpack the given byte array. Since the length
// of the byte array does not contain enough information for the number of
// bits, an extra number is passed to specify the total number of bits.
func NewBitUnpacker(data []byte, nbit uint32) *BitUnpacker {
	if int(nbit) > len(data)*8 {
		return nil
	}
	return &BitUnpacker{
		data, 0, 0, nbit,
	}
}

// Write bits into a bit packer. The given integer are stored in big endian,
// i.e. the higher bits are stored in earlier bytes.
// For the last byte, if it is not full, the bits are stored in higher part of
// this byte.
func (packer *BitPacker) WriteBits(code uint64, nbit uint32) error {
	if int(packer.nbyte) >= MAX_PACKER_SIZE {
		return fmt.Errorf("Packer full!")
	}
	for nbit > 0 {
		if int(packer.nbyte) >= MAX_PACKER_SIZE {
			return fmt.Errorf("Packer full!")
		}
		left := 8 - packer.nbit
		if nbit < left {
			packer.data[packer.nbyte] |= byte(code&uint64(uint32(1<<nbit)-1)) << (left - nbit)
			packer.nbit += nbit
			return nil
		} else {
			packer.data[packer.nbyte] |= byte((code >> (nbit - left)) & uint64(uint32(1<<left)-1))
			packer.nbit = 0
			packer.nbyte += 1
			nbit -= left
			code &= uint64((1 << nbit) - 1)
		}
	}
	return nil
}

// Returns the number of bits that is left in the bit-unpacker.
func (unpacker *BitUnpacker) Left() uint32 {
	return unpacker.size - (unpacker.pbyte*8 + unpacker.pbit)
}

// Read a specific number of bits in the bit unpacker, and store the result
// in the lower bits of an uint64.
func (unpacker *BitUnpacker) ReadBits(nbit uint32) (uint64, error) {
	if unpacker.Left() < nbit {
		return 0, fmt.Errorf("Not enough bits left!")
	}
	ret := uint64(0)
	for nbit > 0 {
		left := 8 - unpacker.pbit
		if nbit < left {
			ret <<= nbit
			ret |= uint64((unpacker.data[unpacker.pbyte] >> (left - nbit)) & byte((1<<nbit)-1))
			unpacker.pbit += nbit
			return ret, nil
		} else {
			ret <<= left
			ret |= uint64(unpacker.data[unpacker.pbyte] & byte((1<<left)-1))
			unpacker.pbyte += 1
			unpacker.pbit = 0
			nbit -= left
		}
	}
	return ret, nil
}

// Returns the number of bits stored in the bit packer.
func (packer *BitPacker) Size() uint32 {
	return packer.nbyte*8 + packer.nbit
}

// Returns the slice of bytes containing the stored bits.
func (packer *BitPacker) Data() []byte {
	return packer.data[:(packer.Size()+7)/8]
}
