package main

import (
	"fmt"

	"github.com/DeanXu2357/go-aiy-voice/aiy/board"
)

func main() {
	fmt.Println("Test Start")

	// ShinePin25()
	board.NewButton(23, 0.25, func() {
		fmt.Println("pressed\n")
	})
}
