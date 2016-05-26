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
  sc 	subscribe
  ssc	show subscribe count
  sm	show/hide received messages
  stm	start count messages
  edm	end message count
  smc   show message count
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
			fmt.Printf("---Connected: %d/%d\n", ConnectionCount, b.C)
		case "sc":
			for i := 0; i < b.C; i++ {
				go b.subscribe(i, "yy1")
			}
			fmt.Print("Subscribe done\n")
		case "ssc":
			fmt.Printf("---SubscribeCount: %d\n", SubscribeCount)
		case "sm":
			ShowReceivedMessages = !ShowReceivedMessages
			fmt.Printf("---ShowReceivedMessages: %v\n", ShowReceivedMessages)
		case "stm":
			MessageCount = 0
			StartCountMessages = true
			//fmt.Printf("\nMessageCount: %v\n", MessageCount)
		case "edm":
			StartCountMessages = false
			fmt.Printf("---MessageCount: %v\n", MessageCount)
		case "smc":
			fmt.Printf("---MessageCount: %v\n\n", MessageCount)
		case "h":
			fallthrough
		case "?":
			fmt.Fprint(os.Stderr, fmt.Sprintf(commands))
		default:
			fmt.Println("unknow command")
		}
	}
}
