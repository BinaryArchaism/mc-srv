package datatypes

import "io"

type Boolean bool

func WriteBoolean(b Boolean) byte {
	if b {
		return 0x01
	}

	return 0x00
}

func ReadBoolean(b byte) Boolean {
	return b == 0x01
}

func (b Boolean) Write(w io.ByteWriter) error {
	err := w.WriteByte(WriteBoolean(b))
	if err != nil {
		return err
	}
	return nil
}

func (b Boolean) Read(r io.ByteReader) (Boolean, error) {
	byteBool, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	return ReadBoolean(byteBool), nil
}
