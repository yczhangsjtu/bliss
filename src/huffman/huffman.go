package huffman

import (
	"fmt"
)

type Pair struct {
	Code uint64
	Nbit uint8
}

type HuffmanCode struct {
	Code []Pair
}

type HuffmanEncoder struct {
	packer *BitPacker
	code   *HuffmanCode
}

type HuffmanDecoder struct {
	unpacker *BitUnpacker
	code     *HuffmanCode
}

func NewHuffmanEncoder(code *HuffmanCode) *HuffmanEncoder {
	return &HuffmanEncoder{NewBitPacker(), code}
}

func NewHuffmanDecoder(code *HuffmanCode, data []byte) *HuffmanDecoder {
	size := uint32(data[0])*256 + uint32(data[1])
	return &HuffmanDecoder{NewBitUnpacker(data[2:], size), code}
}

func (code *HuffmanCode) Size() int {
	return len(code.Code)
}

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

func (encoder *HuffmanEncoder) Digest() []byte {
	size := encoder.packer.Size()
	ret := []byte{byte(size / 256), byte(size % 256)}
	ret = append(ret, encoder.packer.Data()...)
	return ret
}

func (decoder *HuffmanDecoder) Next() int {
	return 0
}
