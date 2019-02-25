package board

import (
	"bufio"
	"fmt"
	"os"
)

type KeyListener struct {
	listenKey string
	scan      chan string
}

func NewKeyListener(listenKey string) (listener listener, err error) {
	var scan chan string

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input := scanner.Text()
			scan <- input
			fmt.Println("Scan : " + input)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	return KeyListener{listenKey: listenKey, scan: scan}, err
}

func (listener KeyListener) IsTriggered() (result bool) {
	result = false

	select {
	case x := <-listener.scan:
		if x == listener.listenKey {
			return true
		}
	default:
		return
	}

	return
}

func (listener KeyListener) Close() {
	//
	close(listener.scan)
}
