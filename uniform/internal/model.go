package internal

const Delimiter = 0x50554646

type TLV struct {
	TLVHeader
	Value []byte
}

type TLVHeader struct {
	Type   uint16
	Length uint32
}

const (
	TLVTypeHeader = iota + 1
	TLVTypeUserData
)
