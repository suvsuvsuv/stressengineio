package boomer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var commands = `
Commands:
  h  	help
  q  	quit
  gc	gc connected count
  gm	get message count
  sc <>	subscribe
`

const inputdelimiter = '\n'

func (b *Boomer) readConsole() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter command: ")
		input, err := reader.ReadString(inputdelimiter)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)
		switch input {
		case "q":
			close(requestCountChan)
			os.Exit(1)
		case "gc":
			fmt.Printf("Connected: %d/%d\n", ConnectionCount, b.C)
		case "h":
			fallthrough
		case "?":
			fmt.Fprint(os.Stderr, fmt.Sprintf(commands))
		default:
			fmt.Println("unknow command")
		}
	}
}
