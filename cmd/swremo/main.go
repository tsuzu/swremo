package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/tsuzu/swremo/pkg/gamecon"
)

const script = `#!/bin/bash

cd /sys/kernel/config/usb_gadget/

mkdir -p gamepad
cd gamepad

echo 0x0f0d > idVendor # Linux Foundation
echo 0x00c1 > idProduct # Multifunction Composite Gadget
echo 0x0572 > bcdDevice # v5.7.2
echo 0x0200 > bcdUSB # USB2
mkdir -p strings/0x409
echo "" > strings/0x409/serialnumber
echo "HORI CO.,LTD." > strings/0x409/manufacturer
echo "HORIPAD S" > strings/0x409/product
mkdir -p configs/c.1/strings/0x409
echo "" > configs/c.1/strings/0x409/configuration
echo 500 > configs/c.1/MaxPower
echo 0x80 configs/c.1/bmAttributes

mkdir -p  functions/hid.usb0
echo 0 >  functions/hid.usb0/protocol
echo 0 >  functions/hid.usb0/subclass
echo 64 > functions/hid.usb0/report_length
echo -ne \\x5\\x1\\x9\\x5\\xa1\\x1\\x15\\x0\\x25\\x1\\x35\\x0\\x45\\x1\\x75\\x1\\x95\\xe\\x5\\x9\\x19\\x1\\x29\\xb\\x81\\x2\\x95\\x2\\x81\\x1\\x5\\x1\\x25\\x7\\x46\\x3b\\x1\\x75\\x4\\x95\\x1\\x65\\x14\\x9\\x39\\x81\\x42\\x65\\x0\\x95\\x1\\x81\\x1\\x26\\xff\\x0\\x46\\xff\\x0\\x9\\x30\\x9\\x31\\x9\\x32\\x9\\x35\\x75\\x8\\x95\\x4\\x81\\x2\\x75\\x8\\x95\\x1\\x81\\x1\\xc0 > functions/hid.usb0/report_desc

ln -s functions/hid.usb0 configs/c.1/

# End functions

udevadm settle -t 5 || :
ls /sys/class/udc > UDC
`

func createGamepad() error {
	return nil

	cmd := exec.Command("bash", "/dev/stdin")
	cmd.Stdin = strings.NewReader(script)

	output, err := cmd.CombinedOutput()

	log.Println(string(output))

	return err
}

type PadJSON struct {
	gamecon.GameController
}

func main() {
	ch := make(chan gamecon.GameController, 1)

	go func(i int) {
		var fp *os.File

		reopen := func() {
			if err := createGamepad(); err != nil {
				log.Println(err)
			}

			var err error
			fp, err = os.OpenFile(fmt.Sprintf("/dev/hidg%d", i), os.O_WRONLY, 0600)

			if err != nil {
				log.Println(err)
			}
		}

		reopen()

		defer func() {
			if fp != nil {
				fp.Close()
			}
		}()

		for {
			data, ok := <-ch

			if !ok {
				return
			}

			_, err := fp.Write(data.HIDRepresentation())

			if err != nil {
				log.Println(err)

				reopen()
			}
		}
	}(0)

	listener, err := net.Listen("tcp", ":61345")

	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println(err)
			continue
		}

		dec := json.NewDecoder(conn)

		for {
			var req PadJSON
			if err := dec.Decode(&req); err != nil {
				conn.Close()
				log.Println(err)
				break
			}

			select {
			case ch <- req.GameController:
			default:
			}
		}
	}

}
