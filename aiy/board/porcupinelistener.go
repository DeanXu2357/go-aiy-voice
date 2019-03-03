package board

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	porcupine "github.com/charithe/porcupine-go"
)

// PorcupineListener , keyword detection listener
type PorcupineListener struct {
	read       chan bool
	done       chan bool
	modelPath  string
	keywords   keywordFlags
	audioInput io.Reader
}

type keywordFlags []*porcupine.Keyword

// NewKeyword , new keyword object
func NewKeyword(value string, filePath string, sensitivity float32) *porcupine.Keyword {
	return &porcupine.Keyword{Value: value, FilePath: filePath, Sensitivity: sensitivity}
}

// NewPorcupineListener , new keyword detection listener
func NewPorcupineListener(audioInput io.Reader, modelPath string, keywords ...*porcupine.Keyword) (listener PorcupineListener, err error) {
	listener = PorcupineListener{
		audioInput: audioInput,
		modelPath:  modelPath,
		keywords:   keywords,
		read:       make(chan bool),
		done:       make(chan bool),
	}

	go listener.listen()

	return
}

func (listener *PorcupineListener) listen() {
	p, err := porcupine.New(listener.modelPath, listener.keywords...)
	if err != nil {
		fmt.Println(err)
		// log
		// log.Fatalf("failed to initialize porcupine: %+v", err)
		os.Exit(2)
		return
	}
	defer p.Close()

	frameSize := porcupine.FrameLength()
	audioFrame := make([]int16, frameSize)
	buffer := make([]byte, frameSize*2)

	fmt.Printf("listening...\n")

	for {
		// fmt.Println("..")

		select {
		case done := <-listener.done:
			if done {
				return
			}
		default:
			// fmt.Println(".")
			if err := readAudioFrame(listener.audioInput, buffer, audioFrame); err != nil {
				fmt.Printf("error: %+v", err)
				// log
				return
			}
			// fmt.Println("readaudio frame success")

			word, err := p.Process(audioFrame)
			if err != nil {
				fmt.Printf("error: %+v", err)
				// log
				continue
			}
			// fmt.Println("porcupine porcess")

			if word != "" {
				fmt.Printf("detected word: \"%s\"", word)
				listener.read <- true
				// log
			}
			// fmt.Println("porcupine detected")
		}
	}
}

func readAudioFrame(src io.Reader, buffer []byte, audioFrame []int16) error {
	_, err := io.ReadFull(src, buffer)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(buffer)
	for i := 0; i < len(audioFrame); i++ {
		if err := binary.Read(buf, binary.LittleEndian, &audioFrame[i]); err != nil {
			return err
		}
	}

	return nil
}

// IsTriggered , if keyword is detected
func (listener PorcupineListener) IsTriggered() bool {
	select {
	case read := <-listener.read:
		return read
	default:
		return false
	}
}

// Close , close listener
func (listener PorcupineListener) Close() {
	listener.done <- true
}

// GetAudioInput , get audio from file path or from stdin default
// To use this method,  has to execute
// ```gst-launch-1.0 -v alsasrc ! audioconvert ! audioresample ! audio/x-raw,channels=1,rate=16000,format=S16LE ! filesink location=/dev/stdout |```
// before run ``` go run main.go ```
func GetAudioInput(filePath string) (audio io.Reader, err error) {
	if filePath == "" {
		audio = bufio.NewReader(os.Stdin)
	} else {
		f, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("failed to open input [%s]: %+v\n", filePath, err)
			os.Exit(2)
		}
		defer f.Close()

		audio = bufio.NewReader(f)
	}

	return
}
