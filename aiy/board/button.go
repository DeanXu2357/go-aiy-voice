package board

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// Button 按鈕
type Button struct {
	pin          uint8
	debounceTime float64
	channel      chan buttonCh
	whenPressed  func()
	whenReleased func()
}

type buttonCh struct {
	done         bool
	status       uint8
	whenPressed  func()
	whenReleased func()
}

const (
	// Released 放開按鈕狀態常數 0
	Released uint8 = iota
	Pressed
)

// NewButton 新增按鈕監聽物件
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

	return
}

func buttonRun(bCh chan buttonCh, button Button) {
	pin := rpio.Pin(button.pin)
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()
	fmt.Println("button run start")

	var ch buttonCh
	pressed := false
	whenPress := time.Now()
	whp := button.whenPressed
	whr := button.whenReleased

	loop := true
	for loop {
		select {
		case ch = <-bCh:
			if ch.done {
				loop = false
			}
		default:
		}

		if time.Since(whenPress).Seconds() > button.debounceTime {
			if pin.Read() == rpio.Low {
				if !pressed {
					pressed = true
					whenPress = time.Now()
					if ch.whenPressed != nil {
						whp = ch.whenPressed
					}
					go whp()
					bCh <- buttonCh{status: Pressed}
				}
			} else {
				if pressed {
					pressed = false
					whenPress = time.Now()
					if ch.whenReleased != nil {
						whr = ch.whenReleased
					}
					go whr()
					bCh <- buttonCh{status: Released}
				}
			}
		}

		if !loop {
			fmt.Println("close channel")
			close(bCh)
		}
	}

	fmt.Println("button run end")
	rpio.Close()
}

// Close 關閉按鈕監聽
func (btn *Button) Close() {
	btn.channel <- buttonCh{done: true}
}

// SetWhenPressed 設定新按鈕按下事件
func (btn *Button) SetWhenPressed(f func()) {
	btn.whenPressed = f
	btn.channel <- buttonCh{whenPressed: f}
}

// SetWhenReleased 設定新按鈕放開事件
func (btn *Button) SetWhenReleased(f func()) {
	btn.whenReleased = f
	btn.channel <- buttonCh{whenReleased: f}
}

// WaitForPressed 暫時凍結等待按鈕按下
func (btn *Button) WaitForPressed() {
	var ch buttonCh
	for {
		select {
		case ch = <-btn.channel:
			if ch.status == Pressed {
				return
			}
		default:
		}

		time.Sleep(1 * time.Second)
	}
}

// WaitForReleased 暫時凍結等待按鈕放開
func (btn *Button) WaitForReleased() {
	var ch buttonCh
	for {
		select {
		case ch = <-btn.channel:
			if ch.status == Released {
				return
			}
		default:
		}

		time.Sleep(1 * time.Second)
	}
}
