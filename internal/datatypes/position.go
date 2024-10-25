package datatypes

type Position struct {
	X int
	Z int
	Y int
}

func ReadPosition(val int64) Position {
	x := int(val >> 38)
	y := int(val << 52 >> 52)
	z := int(val << 26 >> 38)
	return Position{
		X: x,
		Z: z,
		Y: y,
	}
}

func WritePosition(val Position) int64 {
	return (int64(val.X&0x3FFFFFF) << 38) |
		int64((val.Z&0x3FFFFFF)<<12) |
		int64(val.Y&0xFFF)
}
