package main

import (
	"fmt"
	"os"

	"github.com/stianeikeland/go-rpio/v4"
)

func recorder() {
	exec := true

	pin := rpio.Pin(23)

	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	for exec {
		pin.PullUp()
		fmt.Printf("PullUp: %d\n", pin.Read())
		ShinePin25()

		// Pull down and read value
		pin.PullDown()
		fmt.Printf("PullDown: %d\n", pin.Read())
	}
}
