package network

type LoginLayer struct {
	UserSize uint8
	User     []byte
	PassSize uint8
	Pass     []byte
	Payload  []byte
}

func (l *LoginLayer) Type() string {
	return "LoginLayer"
}

func (l *LoginLayer) FromFrame(xb []byte) error {
	b := xb[0]
	xb = xb[1:]
	l.UserSize = b
	if len(xb) >= int(b) {
		return errNotEnoughDataInFrames
	}
	l.User = make([]byte, b)
	for i := 0; i < int(b); i++ {
		l.User[i] = xb[i]
	}
	xb = xb[b:]

	b = xb[0]
	xb = xb[1:]
	l.PassSize = b
	if len(xb) >= int(b) {
		return errNotEnoughDataInFrames
	}
	l.Pass = make([]byte, b)
	for i := 0; i < int(b); i++ {
		l.Pass[i] = xb[i]
	}
	xb = xb[b:]

	l.Payload = xb
	return nil
}

func (l *LoginLayer) ToFrame() ([]byte, error) {
	out := make([]byte, 0)
	out = append(out, l.UserSize)
	out = append(out, l.User...)
	out = append(out, l.PassSize)
	out = append(out, l.Pass...)
	out = append(out, l.Payload...)
	return out, nil
}
