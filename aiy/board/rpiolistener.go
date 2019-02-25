package board

import (
	"fmt"

	"github.com/stianeikeland/go-rpio/v4"
)

type RpioListener struct {
	pin  uint8
	read chan bool
	done chan bool
}

func NewRpioListener(pin uint8) (listener RpioListener, err error) {

	ch := make(chan bool)
	listener = RpioListener{pin: pin, read: ch}

	go rpioRun(listener)

	return
}

func rpioRun(listener RpioListener) {
	pinAction := rpio.Pin(listener.pin)
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		return
	}
	defer rpio.Close()

	for true {
		status := pinAction.Read()
		var read bool
		if status == rpio.Low {
			read = true
		} else {
			read = false
		}
		listener.read <- read

		select {
		case done := <-listener.done:
			if done {
				close(listener.read)
				close(listener.done)

				return
			}
		default:
		}
	}
}

func (listener RpioListener) IsTriggered() bool {
	select {
	case read := <-listener.read:
		return read
	default:
		return false
	}
}

func (listener RpioListener) Close() {
	listener.done <- true
	rpio.Close()
}
