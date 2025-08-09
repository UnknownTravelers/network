package network

import (
	"fmt"
	"math"
)

type Trame interface {
	FromFrame(xb []byte) error
	ToFrame() ([]byte, error)
	AddBytes(b ...byte)
}

type Basic struct {
	// 1 1bit header
	Version uint8 // version 7bit
	Command uint8 // action 8bit
	Payload []byte
}

func (basic *Basic) FromFrame(xb []byte) error {
	b := xb[0]
	if b&0x80 == 0 {
		return fmt.Errorf("invalid frame format: no leading 1")
	}
	xb = xb[1:]
	basic.Version = b & 0x7f

	b = xb[0]
	xb = xb[1:]
	basic.Command = b

	basic.Payload = xb
	return nil
}

func (basic *Basic) ToFrame() ([]byte, error) {
	out := make([]byte, 0)
	out = append(out, 0x80|basic.Version)
	out = append(out, basic.Command)
	out = append(out, basic.Payload...)
	return out, nil
}

func (basic *Basic) AddBytes(b ...byte) {
	basic.Payload = append(basic.Payload, b...)
}

func (basic *Basic) AddUint16(v uint16) {
	b1 := byte(v >> 8)
	b2 := byte(v & 0x0f)
	basic.AddBytes(b1, b2)
}

func (basic *Basic) AddUint32(v uint32) {
	b1 := byte(v >> 24)
	b2 := byte(v >> 16)
	b3 := byte(v >> 8)
	b4 := byte(v & 0x0f)
	basic.AddBytes(b1, b2, b3, b4)
}

func (basic *Basic) AddUint64(v uint64) {
	b1 := byte(v >> 56)
	b2 := byte(v >> 48)
	b3 := byte(v >> 40)
	b4 := byte(v >> 32)
	b5 := byte(v >> 24)
	b6 := byte(v >> 16)
	b7 := byte(v >> 8)
	b8 := byte(v)
	basic.AddBytes(b1, b2, b3, b4, b5, b6, b7, b8)
}

func (basic *Basic) AddFloat32(f float32) {
	v := math.Float32bits(f)
	basic.AddUint32(v)
}

func (basic *Basic) AddFloat64(f float64) {
	v := math.Float64bits(f)
	basic.AddUint64(v)
}
