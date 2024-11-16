package datatypes

import "io"

type Byte byte

func (b *Byte) Write(w io.ByteWriter) error {
	return w.WriteByte(byte(*b))
}

func (b *Byte) Read(r io.ByteReader) error {
	v, err := r.ReadByte()
	if err != nil {
		return err
	}
	*b = Byte(v)
	return nil
}
