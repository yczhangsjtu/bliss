// Package huffman implements the huffman encoding and decoding scheme.
package huffman

import (
	"fmt"
)

// A Pair specifies a consecutive number of bits. Code is used to store the
// bits, and Nbit specifies the number of bits.
type Pair struct {
	Code uint64
	Nbit uint8
}

// A Triple is a tuple (left, right, index) in the decoding graph. Either
// index = -1, or left = right = -1, and they cannot both be -1.
// left is the index of the next triple if a 0 is encountered, and right is
// when a 1 is encountered. If index is not zero, then index is the decoded
// result.
type Triple struct {
	Left  int
	Right int
	Index int
}

// A HuffmanCode is a precomputed huffman encoding and decoding table.
// Code is used to encode, and Node is used to decode.
// The table can encode integers from 0 to len(Code)-1.
// An integer 0 <= i < len(Code) is encoded into the bits specified by Code[i].
type HuffmanCode struct {
	Code []Pair
	Node []Triple
}

// A HuffmanEncoder is a live instance of a huffman encoding procedure.
// It refers to a HuffmanCode for the huffman table, and uses a bit packer to
// store the encoded bits.
type HuffmanEncoder struct {
	packer *BitPacker
	code   *HuffmanCode
}

// A HuffmanEncoder is a live instance of a huffman decoding procedure.
// It refers to a HuffmanCode for the huffman table, and read the bits to be
// decoded from a bit unpacker.
type HuffmanDecoder struct {
	unpacker *BitUnpacker
	code     *HuffmanCode
}

// Create a HuffmanEncoder for a HuffmanCode.
func NewHuffmanEncoder(code *HuffmanCode) *HuffmanEncoder {
	return &HuffmanEncoder{NewBitPacker(), code}
}

// Create a HuffmanDecoder for a HuffmanCode and data to decode.
// The data is in the following format:  | bit-size (2 bytes) | content |.
// The decoder reads the bit size from the first two bytes, and put the content
// in a newly created bit unpacker.
func NewHuffmanDecoder(code *HuffmanCode, data []byte) *HuffmanDecoder {
	size := uint32(data[0])*256 + uint32(data[1])
	return &HuffmanDecoder{NewBitUnpacker(data[2:], size), code}
}

// Returns the size of a huffman code, which is the number of symbols this code
// encodes.
func (code *HuffmanCode) Size() int {
	return len(code.Code)
}

// Encode a symbol represented by index in [0, size).
func (encoder *HuffmanEncoder) Update(index int) error {
	code := encoder.code
	if index < 0 || index >= code.Size() {
		return fmt.Errorf("Invalid symbol %d, expected < %d", index, code.Size())
	}
	err := encoder.packer.WriteBits(code.Code[index].Code, uint32(code.Code[index].Nbit))
	if err != nil {
		return fmt.Errorf("Error in writing bits to packer: %s", err.Error())
	}
	return nil
}

// Return the encoded bits so far. Prepend the number if bits to the data in
// the first two bytes.
func (encoder *HuffmanEncoder) Digest() []byte {
	size := encoder.packer.Size()
	ret := []byte{byte(size / 256), byte(size % 256)}
	ret = append(ret, encoder.packer.Data()[:(size+7)/8]...)
	return ret
}

// Decode the next symbol. This is accomplished by starting from the triple
// Node[curr=0], and read one bit at a time. For a bit of 0, let
// curr = Node[curr]->left, and for a bit of 1, let curr = Node[curr]->right.
// Finally we get somewhere Node[curr]->index is not -1, then we return index
// representing the decoded symbol.
func (decoder *HuffmanDecoder) Next() (int, error) {
	curr := 0
	for decoder.unpacker.Left() > 0 {
		bit, err := decoder.unpacker.ReadBits(1)
		if err != nil {
			return -1, fmt.Errorf("Error in reading bit: %s", err.Error())
		}
		if bit == 0 {
			curr = decoder.code.Node[curr].Left
		} else {
			curr = decoder.code.Node[curr].Right
		}
		if curr < 0 {
			return -1, fmt.Errorf("Unexpected bit %d", bit)
		}
		if decoder.code.Node[curr].Index >= 0 {
			return decoder.code.Node[curr].Index, nil
		}
	}
	return -1, fmt.Errorf("Unexpected end of bit string when decoding")
}
