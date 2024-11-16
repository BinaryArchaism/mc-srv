package datatypes

import (
	"encoding/binary"
	"io"
)

type UShort uint16

func (s *UShort) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, uint16(*s))
}

func (s *UShort) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, s)
}

type Short int16

func (s *Short) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, int16(*s))
}

func (s *Short) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, s)
}
