package datatypes

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
