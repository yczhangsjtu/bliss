package huffman

import (
	"testing"
)

func TestWriteBits(t *testing.T) {
	packer := NewBitPacker()
	packer.WriteBits(0x1, 1)
	packer.WriteBits(0x2, 2)
	packer.WriteBits(0x3, 2)
	packer.WriteBits(0x4, 3)
	packer.WriteBits(0x5, 3)
	packer.WriteBits(0x6, 3)
	expect := []byte{0xdc, 0x2e}
	// Now the bits should be 11011100 101110
	// nbytes = 1, nbit = 6
	if packer.nbyte != 1 {
		t.Errorf("Wrong number of bytes, expected %d, got %d", 1, packer.nbyte)
	}
	if packer.nbit != 6 {
		t.Errorf("Wrong number of bits, expected %d, got %d", 6, packer.nbit)
	}
	for i := 0; i <= int(packer.nbyte); i++ {
		if packer.data[i] != expect[i] {
			t.Errorf("Wrong byte at %d: expected %x, got %x", i, expect[i], packer.data[i])
		}
	}
}
