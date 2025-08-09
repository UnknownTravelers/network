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
			name:    "No leading 1",
			xb:      []byte{0x00},
			wantErr: true,
		},
		{
			name: "No Payload",
			xb:   []byte{0x80, 0x00},
			want: fields{
				Version: 0,
				Command: 0,
				Payload: []byte{},
			},
			wantErr: false,
		},
		{
			name: "Version",
			xb:   []byte{0x81, 0x00},
			want: fields{
				Version: 1,
				Command: 0,
				Payload: []byte{},
			},
			wantErr: false,
		},
		{
			name: "Command",
			xb:   []byte{0x80, 0x02},
			want: fields{
				Version: 0,
				Command: 2,
				Payload: []byte{},
			},
			wantErr: false,
		},
		{
			name: "Payload",
			xb:   []byte{0x80, 0x00, 0x00, 0x01, 0xff},
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
