package main

import (
	"fmt"
	"os"

	"go-aiy-voice/aiy/board"
)

func main() {
	fmt.Println("Test Start")

	// ShinePin25()
	btn, err := board.NewButton(23, 0.25, func() {
		fmt.Println("starter set pressed event\n")
	}, func() {})
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
