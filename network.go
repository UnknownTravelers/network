package network

import (
	"fmt"
	"math"
)

var errNotEnoughData = fmt.Errorf("not enough data in payload")

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
	b2 := byte(v)
	basic.AddBytes(b1, b2)
}

func (basic *Basic) AddUint32(v uint32) {
	b1 := byte(v >> 24)
	b2 := byte(v >> 16)
	b3 := byte(v >> 8)
	b4 := byte(v)
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

func (basic *Basic) PopByte() (byte, error) {
	if len(basic.Payload) <= 0 {
		return 0, errNotEnoughData
	}
	v := basic.Payload[0]
	basic.Payload = basic.Payload[1:]
	return v, nil
}

func (basic *Basic) PopUint16() (uint16, error) {
	if len(basic.Payload) < 2 {
		return 0, errNotEnoughData
	}
	v := basic.Payload[0:2]
	basic.Payload = basic.Payload[2:]
	var out uint16 = (uint16(v[0]) << 8) + uint16(v[1])
	return out, nil
}

func (basic *Basic) PopUint32() (uint32, error) {
	if len(basic.Payload) < 4 {
		return 0, errNotEnoughData
	}
	v := basic.Payload[0:4]
	basic.Payload = basic.Payload[4:]
	var out uint32 = (uint32(v[0]) << 24) + (uint32(v[1]) << 16) + (uint32(v[2]) << 8) + uint32(v[3])
	return out, nil
}

func (basic *Basic) PopUint64() (uint64, error) {
	if len(basic.Payload) < 8 {
		return 0, errNotEnoughData
	}
	v := basic.Payload[0:8]
	basic.Payload = basic.Payload[8:]
	var out uint64 = (uint64(v[0]) << 56) + (uint64(v[1]) << 48) + (uint64(v[2]) << 40) + (uint64(v[3]) << 32) + (uint64(v[4]) << 24) + (uint64(v[5]) << 16) + (uint64(v[6]) << 8) + uint64(v[7])
	return out, nil
}
