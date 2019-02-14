package voice

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Speacker tts parameter bag struct
type Speacker struct {
	Text   string
	Lang   string
	Volume uint8
	Pitch  uint8
	Speed  uint8
	Device string
}

func (param *Speacker) init() {
	if param.Text == "" {
		param.Text = "hello world"
	}

	if param.Lang == "" {
		param.Lang = "en-US"
	}

	if param.Volume == 0 {
		param.Volume = 80
	}

	if param.Pitch == 0 {
		param.Pitch = 130
	}

	if param.Speed == 0 {
		param.Speed = 100
	}

	if param.Device == "" {
		param.Device = "default"
	}
}

// Say Text to Speach
func Say(param Speacker) (err error) {
	param.init()

	body := fmt.Sprintf("\"<volume level='%d'><pitch level='%d'><speed level='%d'>%s</speed></pitch></volume>\"", param.Volume, param.Pitch, param.Speed, param.Text)

	tmpFile := "tmp.wav"

	cmd := exec.Command("pico2wave",
		"--wave="+tmpFile,
		"--lang="+param.Lang,
		body,
		"&& aplay",
		"-q",
		"-D",
		param.Device,
		tmpFile)
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return err
	}

	return
}
