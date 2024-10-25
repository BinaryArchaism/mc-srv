package datatypes

import (
	"encoding/binary"
	"errors"
)

type Short int16

func WriteShort(v Short) []byte {
	res := make([]byte, 2)
	binary.BigEndian.PutUint16(res, uint16(v))
	return res
}

func ReadShort(v []byte) Short {
	if len(v) != 2 {
		panic(errors.New("cannot read uint16"))
	}
	return Short(int16(binary.BigEndian.Uint16(v)))
}

type UShort uint16

func WriteUShort(v UShort) []byte {
	res := make([]byte, 2)
	binary.BigEndian.PutUint16(res, uint16(v))
	return res
}

func ReadUShort(v []byte) UShort {
	if len(v) != 2 {
		panic(errors.New("cannot read UShort"))
	}
	return UShort(binary.BigEndian.Uint16(v))
}
