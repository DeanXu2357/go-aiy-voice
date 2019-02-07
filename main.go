package main

import (
	"fmt"

	"go-aiy-voice/aiy/board"
)

func main() {
	fmt.Println("Test Start")

	// ShinePin25()
	board.NewButton(23, 0.25, func() {
		fmt.Println("pressed\n")
	})
}
