package gamecon

import "encoding/binary"

const (
	HATTop uint8 = iota
	HATTopRight
	HATRight
	HATBottomRight
	HATBottom
	HATBottomLeft
	HATLeft
	HATTopLeft
	HATDefault
)

const (
	ButtonY uint = 1 << iota
	ButtonB
	ButtonA
	ButtonX
	ButtonL
	ButtonR
	ButtonZL
	ButtonZR
	ButtonSelect
	ButtonStart
	ButtonLClick
	ButtonRClick
	ButtonHome
	ButtonCapture
)

type GameController struct {
	Button uint
	HAT    uint8
	Axis   [4]uint8
}

func (c *GameController) HIDRepresentation() []byte {
	buf := make([]byte, 8)

	button := c.Button | uint(c.HAT)*(1<<16)

	binary.LittleEndian.PutUint32(buf, uint32(button))
	buf[3] = byte(c.Axis[0]) // Left X
	buf[4] = byte(c.Axis[1]) // Left Y
	buf[5] = byte(c.Axis[2]) // Right X
	buf[6] = byte(c.Axis[3]) // Right Y

	return buf
}
