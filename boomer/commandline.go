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
  shc	show connected count
  shm	show message count
  sc 	subscribe
  ssc	show subscribe count
  sm	show/hide received messages
  stm	start counting messages
  edm	end counting message count
`

const inputdelimiter = '\n'

var doneSuscribe bool

func (b *Boomer) readConsole() {
	doneSuscribe = false
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
		case "shc":
			fmt.Printf("---Connected: %d/%d\n", ConnectionCount, b.C)
		case "sc":
			if int(ConnectionCount) != b.C {
				fmt.Printf("---Please wait till all clients are connected: %d/%d\n",
					ConnectionCount, b.C)
				continue
			} else if doneSuscribe {
				fmt.Print("---Already subscribe\n")
				continue
			}
			doneSuscribe = true
			for i := 0; i < b.C; i++ {
				go b.subscribe(i, "yy1")
			}
			fmt.Print("---Subscribing, please wait\n")
		case "ssc":
			fmt.Printf("---SubscribeCount: %d\n", SubscribeCount)
		case "sm":
			ShowReceivedMessages = !ShowReceivedMessages
			fmt.Printf("---ShowReceivedMessages: %v\n", ShowReceivedMessages)
		case "stm":
			MessageCount = 0
			StartCountMessages = true
			fmt.Printf("---MessageCount: %v\n", MessageCount)
		case "edm":
			StartCountMessages = false
			fmt.Printf("---MessageCount: %v\n", MessageCount)
		case "shm":
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
