package main

import (
	"fmt"
	"os"

	"go-aiy-voice/aiy/board"
	// "go-aiy-voice/aiy/voice"
)

func main() {
	fmt.Println("Test Start")
	listener, err := board.NewRpioListener(23)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// ShinePin25()
	btn, err := board.NewButton(listener, 0.25, func() {
		fmt.Println("starter set pressed event")
	}, func() {
		fmt.Println("starter set released event")
		// go voice.Aplay("~/test6.wav", "default")
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		btn.WaitForPressed()
		fmt.Println("test pressed")
		btn.WaitForReleased()
		fmt.Println("test released")
	}
}
