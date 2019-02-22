package voice

import (
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// FilePlayer , file player
type FilePlayer struct {
	doneCh chan bool
}

// NewFilePlayer , new file player
func NewFilePlayer() (player FilePlayer, err error) {
	return
}

// Play , play file
func (f *FilePlayer) Play(fileName string, fileType string, device string) (err error) {
	options, err := producePlayOptions(fileName, fileType, device)
	if err != nil {
		return err
	}

	go playAudio(options, f.doneCh)
	return
}

func playAudio(options []string, endSignal chan bool) {
	cmd := exec.Command("aplay", options...)
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

func producePlayOptions(fileName string, fileType string, device string) (options []string, err error) {
	if !validateFileType(fileType) {
		err = errors.New(errorFileType)
		return
	}

	options = []string{
		"-q",
		"-D", device,
		"-t", fileType,
	}

	if fileName != "" {
		options = append(options, fileName)
	}

	return
}

// Stop , stop playing
func (f *FilePlayer) Stop() {
	f.doneCh <- true
}
