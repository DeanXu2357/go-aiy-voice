package board

import (
	"fmt"

	"github.com/stianeikeland/go-rpio/v4"
)

type RpioListener struct {
	pin rpio.Pin
}

func NewRpioListener(pin uint8) (listener listener, err error) {
	pinAction := rpio.Pin(pin)
	if err = rpio.Open(); err != nil {
		fmt.Println(err)
		// os.Exit(1)
		return
	}
	defer rpio.Close()

	listener = RpioListener{pin: pinAction}

	return
}

func (listener RpioListener) IsTriggered() bool {
	return listener.pin.Read() == rpio.Low
}

func (listener RpioListener) End() {
	rpio.Close()
}
