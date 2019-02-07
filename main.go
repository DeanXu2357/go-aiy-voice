package main

import (
	"fmt"

	"go-aiy-voice/aiy/board"
)

func main() {
	fmt.Println("Test Start")

	// ShinePin25()
	btn := board.NewButton(23, 0.25, func() {
		fmt.Println("pressed\n")
	}, func() {})

	for {
		btn.WaitForPressed()
		fmt.Println("test pressed")
		btn.WaitForReleased()
		fmt.Println("test released")
	}
}
