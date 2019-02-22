package voice

import (
	"fmt"
	"os/exec"
	"time"
)

// Recorder , the record
type Recorder struct {
	startEvent func()
	stopEvent  func()
	doneCh     chan bool
	options    []string
}

// NewRecorder , create a new Recorder struct
func NewRecorder(starter func(), ender func(), afmt AudioFormat) (recorder Recorder, err error) {
	options := produceCmdOptions(afmt)

	endSignal := make(chan bool)

	recorder = Recorder{startEvent: starter, stopEvent: ender, doneCh: endSignal, options: options}

	return
}

func recordRun(options []string, endSignal chan bool) {
	cmd := exec.Command("arecord", options...)
	if err := cmd.Start(); err != nil {
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
	case end := <-endSignal:
		if end {
			if err := cmd.Process.Kill(); err != nil {
				fmt.Printf("failed to kill process %v\n", err)
			}
			fmt.Println("processkilled by end trigger")
			close(endSignal)
		}
	}
}

func produceCmdOptions(afmt AudioFormat) (cmdString []string) {
	return
}

// Start , start recording
func (recorder *Recorder) Start() {
	go recordRun(recorder.options, recorder.doneCh)
}

// Stop , stop recording
func (recorder *Recorder) Stop() {
	recorder.doneCh <- true
}

// SetStartEvent , set the recorder function triggered at start
func (recorder *Recorder) SetStartEvent(fc func()) {
	recorder.startEvent = fc
}

// SetStopEvent , set the recorder function triggered at stop
func (recorder *Recorder) SetStopEvent(fc func()) {
	recorder.stopEvent = fc
}
