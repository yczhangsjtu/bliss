package huffman

type Pair struct {
	Code uint64
	Nbit uint8
}

type HuffmanCode struct {
	Code []Pair
}
