package network

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasic(t *testing.T) {
	type fields struct {
		Version uint8
		Command uint8
		Payload []byte
	}
	tests := []struct {
		name    string
		want    fields
		xb      []byte
		wantErr bool
	}{
		{
			name: "No Payload",
			xb:   []byte{0x00, 0x00},
			want: fields{
				Version: 0,
				Command: 0,
				Payload: []byte{},
			},
			wantErr: false,
		},
		{
			name: "Version",
			xb:   []byte{0x01, 0x00},
			want: fields{
				Version: 1,
				Command: 0,
				Payload: []byte{},
			},
			wantErr: false,
		},
		{
			name: "Command",
			xb:   []byte{0x00, 0x02},
			want: fields{
				Version: 0,
				Command: 2,
				Payload: []byte{},
			},
			wantErr: false,
		},
		{
			name: "Payload",
			xb:   []byte{0x00, 0x00, 0x00, 0x01, 0xff},
			want: fields{
				Version: 0,
				Command: 0,
				Payload: []byte{0x00, 0x01, 0xff},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Basic{}
			err := got.FromFrame(tt.xb)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				expect := &Basic{
					Version: tt.want.Version,
					Command: tt.want.Command,
					Payload: tt.want.Payload,
				}
				require.Equal(t, expect, got)
				out, err := got.ToFrame()
				require.NoError(t, err)
				require.Equal(t, tt.xb, out)
			}
		})
	}
}

func TestBasicReadWrite(t *testing.T) {
	tests := []struct {
		name    string
		do      func(t *testing.T)
		want    any
		wantErr bool
	}{
		{
			name: "Byte OK",
			do: func(t *testing.T) {
				b := &Basic{}
				var want byte = 0x56
				b.PushBytes(want)
				got, err := b.PopByte()
				require.NoError(t, err)
				require.Equal(t, want, got)
			},
		},
		{
			name: "Byte Err",
			do: func(t *testing.T) {
				b := &Basic{}
				_, err := b.PopByte()
				require.Error(t, err)
			},
		},
		{
			name: "Bytes OK",
			do: func(t *testing.T) {
				b := &Basic{}
				var want uint32 = 0x5603f86a
				b.PushUint32(want)
				got, err := b.PopBytes(4)
				require.NoError(t, err)
				require.Equal(t, []byte{0x56, 0x03, 0xf8, 0x6a}, got)
			},
		},
		{
			name: "Bytes Err",
			do: func(t *testing.T) {
				b := &Basic{}
				var want uint16 = 0x5603
				b.PushUint16(want)
				_, err := b.PopBytes(4)
				require.Error(t, err)
			},
		},
		{
			name: "Uint16 OK",
			do: func(t *testing.T) {
				b := &Basic{}
				var want uint16 = 0x5603
				b.PushUint16(want)
				got, err := b.PopUint16()
				require.NoError(t, err)
				require.Equal(t, want, got)
			},
		},
		{
			name: "Uint16 Err",
			do: func(t *testing.T) {
				b := &Basic{}
				_, err := b.PopUint16()
				require.Error(t, err)
			},
		},
		{
			name: "Uint32 OK",
			do: func(t *testing.T) {
				b := &Basic{}
				var want uint32 = 0x5603f86a
				b.PushUint32(want)
				got, err := b.PopUint32()
				require.NoError(t, err)
				require.Equal(t, want, got)
			},
		},
		{
			name: "Uint32 Err",
			do: func(t *testing.T) {
				b := &Basic{}
				_, err := b.PopUint32()
				require.Error(t, err)
			},
		},
		{
			name: "Uint64 OK",
			do: func(t *testing.T) {
				b := &Basic{}
				var want uint64 = 0x5603f86a5603f86a
				b.PushUint64(want)
				got, err := b.PopUint64()
				require.NoError(t, err)
				require.Equal(t, want, got)
			},
		},
		{
			name: "Uint64 in, Uint32+Uint32 out",
			do: func(t *testing.T) {
				b := &Basic{}
				var want uint32 = 0x5603f86a
				b.PushUint64((uint64(want) << 32) + uint64(want))

				got, err := b.PopUint32()
				require.NoError(t, err)
				require.Equal(t, want, got)

				got, err = b.PopUint32()
				require.NoError(t, err)
				require.Equal(t, want, got)

				_, err = b.PopUint32()
				require.Error(t, err)
			},
		},
		{
			name: "Uint64 Fail",
			do: func(t *testing.T) {
				b := &Basic{}
				_, err := b.PopUint64()
				require.Error(t, err)
			},
		},
		{
			name: "Float32 OK",
			do: func(t *testing.T) {
				b := &Basic{}
				var want float32 = 7.265894e25
				b.PushFloat32(want)
				got, err := b.PopFloat32()
				require.NoError(t, err)
				require.Equal(t, want, got)
			},
		},
		{
			name: "Float64 OK",
			do: func(t *testing.T) {
				b := &Basic{}
				var want float64 = 7.265894e25
				b.PushFloat64(want)
				got, err := b.PopFloat64()
				require.NoError(t, err)
				require.Equal(t, want, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.do)
	}
}
