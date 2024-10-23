package datatypes

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToLittleEndian(t *testing.T) {
	in := []byte{0x00, 0x00, 0x00, 0x0A}
	res := ToLittleEndian(in)
	t.Logf("res: %v", res)
}

func TestVarInt_Bytes(t *testing.T) {
	testCases := []struct {
		inVarInt VarInt
		expOut   []byte
	}{
		{
			inVarInt: 0,
			expOut:   []byte{0x00},
		}, {
			inVarInt: 1,
			expOut:   []byte{0x01},
		}, {
			inVarInt: 128,
			expOut:   []byte{0x80, 0x01},
		}, {
			inVarInt: 255,
			expOut:   []byte{0xff, 0x01},
		}, {
			inVarInt: -12345,
			expOut:   []byte{0xc7, 0x9f, 0x7f},
		}, {
			inVarInt: -1,
			expOut:   []byte{0xff, 0xff, 0xff, 0xff, 0x0f},
		}, {
			inVarInt: -2147483648,
			expOut:   []byte{0x80, 0x80, 0x80, 0x80, 0x08},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d", tc.inVarInt), func(t *testing.T) {
			out := tc.inVarInt.Bytes()
			require.Equal(t, tc.expOut, out)
		})
	}
}

func TestVarInt_Int(t *testing.T) {
	testCases := []struct {
		inBytes []byte
		expOut  VarInt
	}{
		{
			inBytes: []byte{0x00},
			expOut:  0,
		}, {
			inBytes: []byte{0x01},
			expOut:  1,
		}, {
			inBytes: []byte{0x80, 0x01},
			expOut:  128,
		}, {
			inBytes: []byte{0xff, 0x01},
			expOut:  255,
		}, {
			inBytes: []byte{0xff, 0xff, 0xff, 0xff, 0x0f},
			expOut:  -1,
		}, {
			inBytes: []byte{0x80, 0x80, 0x80, 0x80, 0x08},
			expOut:  -2147483648,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%x", tc.inBytes), func(t *testing.T) {
			out, err := ParseVarInt(tc.inBytes)
			require.NoError(t, err)
			require.Equal(t, tc.expOut, out)
		})
	}
}

func TestVarInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		in := rand.Int31n(1000)
		t.Run(fmt.Sprintf("%d", in), func(t *testing.T) {
			out, err := ParseVarInt(VarInt(in).Bytes())
			require.NoError(t, err)
			require.Equal(t, VarInt(in), out)
		})
	}
}
