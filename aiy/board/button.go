package board

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type Button struct {
	pin          uint8
	debounceTime float64
	channel      chan buttonCh
	whenPressed  func()
	whenReleased func()
}

type buttonCh struct {
	done   bool
	status uint8
}

const (
	Released uint8 = iota
	Pressed
)

func NewButton(pin uint8, debounceTime float64, whenPressed func(), whenReleased func()) (button Button, err error) {
	if debounceTime == 0.0 {
		debounceTime = 0.25
	}

	if pin == 0 {
		pin = 23
	}

	bCh := make(chan buttonCh)
	button = Button{
		pin:          pin,
		debounceTime: debounceTime,
		channel:      bCh,
		whenPressed:  whenPressed,
		whenReleased: whenReleased,
	}

	go buttonRun(button.channel, button)
	go doWhenPressed(button)
	go doWhenReleased(button)

	return
}

func buttonRun(bCh chan buttonCh, button Button) {
	whenPress := time.Now()
	pin := rpio.Pin(button.pin)
	pressed := false
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	for true {
		channel := <-bCh
		if channel.done == true {
			break
		}

		if time.Since(whenPress).Seconds() > button.debounceTime {
			if pin.Read() == rpio.Low {
				if !pressed {
					pressed = true
					whenPress = time.Now()
					bCh <- buttonCh{status: Pressed}
				}
			} else {
				if pressed {
					pressed = false
					whenPress = time.Now()
					bCh <- buttonCh{status: Released}
				}
			}
		}
	}

	close(bCh)
	rpio.Close()
}

func doWhenPressed(btn Button) {
	for true {
		if x := <-btn.channel; x.done == true {
			break
		}
		defer close(btn.channel)

		if ch := <-btn.channel; ch.status == Pressed {
			whp := btn.whenPressed
			whp()
		}
	}

	close(btn.channel)
}

func doWhenReleased(btn Button) {
	for true {
		if x := <-btn.channel; x.done == true {
			break
		}
		defer close(btn.channel)

		if ch := <-btn.channel; ch.status == Released {
			whr := btn.whenReleased
			whr()
		}
	}

	close(btn.channel)

}

func (button *Button) Close() {
	button.channel <- buttonCh{done: true}
}

func (btn *Button) SetWhenPressed(f func()) {
	btn.whenPressed = f
}

func (btn *Button) SetWhenReleased(f func()) {
	btn.whenReleased = f
}
