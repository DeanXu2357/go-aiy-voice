package voice

import (
	"fmt"
	"os/exec"
)

func arecord() {

}

func Aplay(filePath string, device string) {
	cmd := exec.Command("aplay", "-D", device, filePath)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %v\n", err)
	}
}
