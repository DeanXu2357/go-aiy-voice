package voice

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Arecord() {

}

func Aplay(filePath string, device string) {
	cmd := exec.Command("aplay", "-D", device, filePath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
}
