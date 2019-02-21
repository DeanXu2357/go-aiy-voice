package voice

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

type VoiceForm struct{}

type Recorder struct {
	startEvent func()
	endEvent   func()
	doneCh     chan bool
}

func NewRecorder(starter func(), ender func()) (recorder Recorder, err error) {
	// cmdString := produceCmdString()

	go func() {
		cmd := exec.Command("aplay", "-D", "", "")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		if err = cmd.Start(); err != nil {
			// todo log
			return
		}

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case <-time.After(3 * time.Second):
			if err := cmd.Process.Kill(); err != nil {
				fmt.Printf("failed to kill process %v\n", err)
			}
			fmt.Println("processkilled as timeout reached")
		case err := <-done:
			if err != nil {
				fmt.Printf("process finished with error = %v", err)
			}
			fmt.Println("process finished successfully")
		}
	}()

	// fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	// fmt.Println("Result: " + out.String())

	return
}

func produceCmdString() string {
	return ""
}

func (recorder *Recorder) Start() {}

func (recorder *Recorder) End() {}

func (recorder *Recorder) IsExist() {}

func (recorder *Recorder) SetStartEvent() {}

func (recorder *Recorder) SetEndEvent() {}
