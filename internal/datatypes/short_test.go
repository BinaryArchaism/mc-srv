package datatypes

import (
	"encoding/binary"
	"github.com/BinaryArchaism/mc-srv/internal/countingbuffer"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"testing"
)

func TestUShort_Read(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		output  UShort
	}{
		{
			name: "zero",
			args: args{
				r: countingbuffer.New(make([]byte, 10)),
			},
			wantErr: false,
			output:  0,
		},
		{
			name: "positive",
			args: args{
				r: countingbuffer.New([]byte{0x00, 0x01, 0x02}),
			},
			wantErr: false,
			output:  1,
		},
		{
			name: "positive2",
			args: args{
				r: countingbuffer.New([]byte{0x10, 0x01, 0x02}),
			},
			wantErr: false,
			output:  0x1001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s UShort
			err := s.Read(tt.args.r)
			if tt.wantErr {
				require.Error(t, err)
			}
			require.NoError(t, err)
			require.Equal(t, tt.output, s)
		})
	}
}

func TestUShort_Write(t *testing.T) {
	tests := []struct {
		name    string
		s       UShort
		wantW   []byte
		wantErr bool
	}{
		{
			name:    "zero",
			s:       0,
			wantW:   []byte{0x00, 0x00},
			wantErr: false,
		},
		{
			name:    "positive",
			s:       0x0123,
			wantW:   []byte{0x01, 0x23},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := countingbuffer.New(make([]byte, 0, 2))
			err := tt.s.Write(buf)
			if tt.wantErr {
				require.Error(t, err)
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantW, buf.Bytes())
		})
	}
}

func TestUShort(t *testing.T) {
	for i := 0; i < 100; i++ {
		rnd := uint16(rand.Int31())
		buf := countingbuffer.New(make([]byte, 0, 2))
		err := binary.Write(buf, binary.BigEndian, rnd)
		require.NoError(t, err)

		var unsignedShort UShort
		err = unsignedShort.Read(buf)
		require.NoError(t, err)
		require.Equal(t, rnd, uint16(unsignedShort))

		expectedBytes := countingbuffer.New(make([]byte, 0, 2))
		err = binary.Write(expectedBytes, binary.BigEndian, rnd)
		require.NoError(t, err)

		buf.Reset()
		err = unsignedShort.Write(buf)
		require.NoError(t, err)
		require.Equal(t, expectedBytes.Bytes(), buf.Bytes())
	}
}

func TestShort(t *testing.T) {
	for i := 0; i < 100; i++ {
		rnd := int16(rand.Int31())
		buf := countingbuffer.New(make([]byte, 0, 2))
		err := binary.Write(buf, binary.BigEndian, rnd)
		require.NoError(t, err)

		var signedShort Short
		err = signedShort.Read(buf)
		require.NoError(t, err)
		require.Equal(t, rnd, int16(signedShort))

		expectedBytes := countingbuffer.New(make([]byte, 0, 2))
		err = binary.Write(expectedBytes, binary.BigEndian, rnd)
		require.NoError(t, err)

		buf.Reset()
		err = signedShort.Write(buf)
		require.NoError(t, err)
		require.Equal(t, expectedBytes.Bytes(), buf.Bytes())
	}
}
