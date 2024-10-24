package datatypes

import (
	"encoding/binary"
	"errors"
)

type Uint16 uint16

func WriteUint16(v Uint16) []byte {
	res := make([]byte, 2)
	binary.BigEndian.PutUint16(res, uint16(v))
	return res
}

func ReadUint16(v []byte) Uint16 {
	if len(v) != 2 {
		panic(errors.New("cannot read uint16"))
	}
	return Uint16(binary.BigEndian.Uint16(v))
}
