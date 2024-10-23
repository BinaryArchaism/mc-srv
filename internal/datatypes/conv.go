package datatypes

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	ErrInvalidVarInt = errors.New("invalid VarInt")
)

func ToLittleEndian(in []byte) int {
	return int(binary.BigEndian.Uint32(in))
}

const (
	segmentBits uint8 = 0x7F // binary 0111 1111
	continueBit uint8 = 0x80 // binary 1000 0000
	shift             = 7
	int32Len          = 32
)

// VarInt for LEB128 representation
type VarInt int32

func (v VarInt) Bytes() []byte {
	const SEGMENT_BITS int32 = 0x7F
	const CONTINUE_BIT int32 = 0x80
	buf := bytes.NewBuffer([]byte{})
	for {
		fmt.Printf("%x\n", v)
		ck := byte(v)
		if (ck & ^(segmentBits)) == 0 {
			buf.WriteByte(byte(v))
			break
		}
		buf.WriteByte(byte((int32(v) & SEGMENT_BITS) | CONTINUE_BIT))
		v >>= shift
	}
	return buf.Bytes()
}

func ParseVarInt(in []byte) (VarInt, error) {
	var (
		res int32
		pos int32
	)

	if len(in) > 5 {
		return 0, ErrInvalidVarInt
	}

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
