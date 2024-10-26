package datatypes

import (
	"encoding/binary"
	"errors"
	"io"
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

func (v *UShort) Read(r io.Reader) error {
	b := make([]byte, 2)
	_, err := r.Read(b)
	if err != nil {
		return err
	}
	*v = UShort(binary.BigEndian.Uint16(b))
	return nil
}

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
