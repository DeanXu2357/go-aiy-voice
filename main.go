package main

import (
	"fmt"
	"go-aiy-voice/aiy/board"
	"os"

	"go-aiy-voice/aiy/voice"

	"github.com/spf13/viper"
)

func initViper() {
	fmt.Println("init viper...")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	initViper()

	fmt.Println("Assistant Start")

	recordFilePath := "/tmp/recorder_tmp.raw"

	afmt := voice.NewAudioFormat(16000, 1, 2)

	recorder, _ := voice.NewRecorder()
	err := recorder.Record(func() {}, func() {}, afmt, "default", recordFilePath)
	if err != nil {
		os.Exit(1)
	}

	// audio, err := board.GetAudioInput(recordFilePath)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	computerKeyword := board.NewKeyword(
		"computer",
		"./resources/porcupine/keyword/computer_linux.ppn",
		0.8,
	)

	porcupineListener, err := board.NewPorcupineListener(
		recorder.Output,
		"./resources/porcupine/lib/porcupine_params.pv",
		computerKeyword)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// fmt.Println("read config")
	// fmt.Println(viper.GetString("MODEL_PATH"))

	for {
		if porcupineListener.IsTriggered() {
			fmt.Println("detected")
			// voice.Say(voice.Speacker{Text: "hello, may i help u ?"})
			// pico2wave --wave=/tmp/test.wav --lang=en-US "<volume level='80'><pitch level='130'><speed level='100'>hello, may i help you?</speed></pitch></volume>" && aplay -q -D default /tmp/test.wav&& rm /tmp/test.wav

		}
	}
}
