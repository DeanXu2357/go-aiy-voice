package voice

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

const (
	errorFileType = "File type must be wav raw au voc"
)

var wg sync.WaitGroup

// Recorder , the record
type Recorder struct {
	doneCh chan bool
	Output io.Reader
}

// NewRecorder , create a new Recorder struct
func NewRecorder() (recorder Recorder, err error) {

	endSignal := make(chan bool)

	recorder = Recorder{doneCh: endSignal}

	return
}

// Record , start recording
func (recorder *Recorder) Record(onStart func(), onStop func(), afmt AudioFormat, device string, fileName string) (err error) {
	options, err := produceCmdOptions(afmt, device, "raw", "")
	if err != nil {
		return err
	}

	go recordRun(options, recorder.doneCh, onStart, onStop, recorder.Output)

	return
}

func recordRun(options []string, endSignal chan bool, onStart func(), onStop func(), output io.Reader) {
	if onStart != nil {
		onStart()
	}

	cmd := exec.Command("arecord", options...)
	if err := cmd.Start(); err != nil {
		// todo log
		panic(err)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	output = bufio.NewReader(&out)
	fmt.Printf("The data is %s\n", out.String())

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
	default:
		// for {
		// }
	}

	if onStop != nil {
		onStop()
	}
}

func produceCmdOptions(afmt AudioFormat, device string, fileType string, fileName string) (cmdString []string, err error) {
	if !validateFileType(fileType) {
		err = errors.New(errorFileType)
		return
	}

	cmdString = []string{
		"-q",
		"-D", device,
		"-t", fileType,
		"-c", strconv.FormatInt(afmt.numChannels, 10),
		"-f", "s" + strconv.FormatInt(afmt.bytesPerSample*8, 10),
		"-r", strconv.FormatInt(afmt.sampleRate, 10),
	}

	if fileName != "" {
		cmdString = append(cmdString, fileName)
	}

	return
}

func validateFileType(fileType string) bool {
	for _, supportType := range supportedFileType {
		if supportType == fileType {
			return true
		}
	}

	return false
}

// Done , stop recording
func (recorder *Recorder) Done() {
	recorder.doneCh <- true
}
