package main

import (
	"fmt"
	"os"

	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	exec := true

	pin := rpio.Pin(23)

	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	pin.Input()

	var res rpio.State

	for exec {
		res = pin.Read()

		if res == rpio.High {
			shinePin25()
		}
	}
}
