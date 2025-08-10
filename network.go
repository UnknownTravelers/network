// Toolbox for network package Creation and Parsing
//
// Made to use with zeromq as communication between game client/server
package network

import (
	"fmt"
)

var errNotEnoughDataInPayload = fmt.Errorf("not enough data in payload")
var errNotEnoughDataInFrames = fmt.Errorf("not enough data in frames")

type Trame interface {
	FromFrame(xb []byte) error
	ToFrame() ([]byte, error)
	Type() string
}
