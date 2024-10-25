package datatypes

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadPosition(t *testing.T) {
	input := int64(0b0100011000000111011000110010110000010101101101001000001100111111)
	require.Equal(t, Position{
		X: 18357644,
		Z: -20882616,
		Y: 831,
	}, ReadPosition(input))
}

func TestWritePosition(t *testing.T) {
	output := int64(0b0100011000000111011000110010110000010101101101001000001100111111)
	pos := Position{
		X: 18357644,
		Z: -20882616,
		Y: 831,
	}
	require.Equal(t, output, WritePosition(pos))
}
