package main

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/tsuzu/joystick"
	"github.com/tsuzu/swremo/pkg/gamecon"
)

func main() {
	gp, err := joystick.Open(0)

	if err != nil {
		panic(err)
	}
	defer gp.Close()

	conn, err := net.Dial("tcp", os.Args[1]+":61345")

	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(conn)

	var prev gamecon.GameController
	for {
		state, err := gp.Read()

		if err != nil {
			panic(err)
		}

		controller := gamecon.GameController{
			Button: uint(state.Buttons),
		}

		xHAT := state.AxisData[4]
		yHAT := state.AxisData[5]
		if xHAT == 0 && yHAT == 0 {
			controller.HAT = gamecon.HATDefault
		} else if xHAT == 0 {
			if yHAT > 0 {
				controller.HAT = gamecon.HATBottom
			} else {
				controller.HAT = gamecon.HATTop
			}
		} else if yHAT == 0 {
			if xHAT > 0 {
				controller.HAT = gamecon.HATRight
			} else {
				controller.HAT = gamecon.HATLeft

			}
		} else if xHAT > 0 {
			switch {
			case yHAT > 0:
				controller.HAT = gamecon.HATBottomRight
			case yHAT == 0:
				controller.HAT = gamecon.HATRight
			default:
				controller.HAT = gamecon.HATTopRight
			}
		} else {
			switch {
			case yHAT > 0:
				controller.HAT = gamecon.HATBottomLeft
			case yHAT == 0:
				controller.HAT = gamecon.HATLeft
			default:
				controller.HAT = gamecon.HATTopLeft
			}
		}

		conv := func(v int) uint8 {
			return uint8((v + 32767) * 255 / 65535)
		}

		controller.Axis[0] = conv(state.AxisData[0])
		controller.Axis[1] = conv(state.AxisData[1])
		controller.Axis[2] = conv(state.AxisData[2])
		controller.Axis[3] = conv(state.AxisData[3])

		if prev == controller {
			continue
		}
		prev = controller

		if err := enc.Encode(controller); err != nil {
			log.Println(err)
			break
		}
	}
}
