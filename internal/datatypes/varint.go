package datatypes

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrInvalidVarInt = errors.New("invalid VarInt")
)

const (
	segmentBits uint8 = 0x7F // binary 0111 1111
	continueBit uint8 = 0x80 // binary 1000 0000
	shift             = 7
	int32Len          = 32
)

// VarInt for LEB128 representation
type VarInt int32

func WriteVarInt(v VarInt) []byte {
	buf := bytes.NewBuffer([]byte{})
	b := uint32(v)
	for {
		if b & ^(uint32(segmentBits)) == 0 {
			buf.WriteByte(byte(b))
			return buf.Bytes()
		}
		buf.WriteByte(byte((b & uint32(segmentBits)) | uint32(continueBit)))
		b = b >> shift
	}
}

func ReadVarInt(in []byte) (VarInt, error) {
	var (
		res int32
		pos int32
	)
	for _, b := range in {
		res |= (int32(b & segmentBits)) << pos
		if b&continueBit == 0 {
			break
		}
		pos += shift
	}

	if pos >= int32Len {
		return 0, ErrInvalidVarInt
	}

	return VarInt(res), nil
}

func ReadVarIntN(in []byte) (VarInt, int, error) {
	var (
		res int32
		pos int32
		n   int
	)
	for _, b := range in {
		n++
		res |= (int32(b & segmentBits)) << pos
		if b&continueBit == 0 {
			break
		}
		pos += shift
	}

	if pos >= int32Len {
		return 0, 0, ErrInvalidVarInt
	}

	return VarInt(res), n, nil
}

func BinaryReadVarInt(in io.ByteReader) (int, error) {
	pl, err := binary.ReadUvarint(in)
	if err != nil {
		return 0, err
	}
	return int(pl), nil
}

func BinaryWriteVarInt(v int) []byte {
	if v < 0 {
		panic(ErrInvalidVarInt)
	}
	b := make([]byte, 5)
	b = b[:binary.PutUvarint(b, uint64(v))]
	return b
}
