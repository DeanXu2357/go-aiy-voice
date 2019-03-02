package main

import (
	"fmt"
	"go-aiy-voice/aiy/board"
	"os"

	// "go-aiy-voice/aiy/voice"

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

	audio, err := board.GetAudioInput("")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	computerKeyword := board.NewKeyword(
		"computer",
		"./resources/porcupine/keyword/computer_linux.ppn",
		0.8,
	)

	porcupineListener, err := board.NewPorcupineListener(
		audio,
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
		}
	}
}
