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
  usc   unsubscribe
  ssc	show subscribe count
  sm	show/hide received messages
  stm	start counting messages
  edm	end counting message count
  tst1  run test 1
`

const inputdelimiter = '\n'

func (b *Boomer) readConsole() {
	doneSuscribe := false
	topicName := "yy1"
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
		ss := strings.Split(input, " ")
		if len(ss) == 0 {
			continue
		}
		switch ss[0] {
		case "q":
			close(requestCountChan)
			os.Exit(1)
		case "shc":
			fmt.Printf("---Connected: %d/%d\n", ConnectionCount, b.C)
		case "sc":
			if len(ss) < 2 {
				fmt.Print("Must provide topicName, e.g.: sc yy1\n")
				continue
			}
			topicName = ss[1]
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
				go func(idx int) {
					b.sendPushID(idx)
					b.subscribe(idx, topicName,
						//"u_live_data:57577d090000654e7c408add"
						nil)
				}(i)
			}
			fmt.Printf("---Subscribing to topic: %s, please wait\n", topicName)
		case "usc":
			if int(ConnectionCount) != b.C {
				fmt.Printf("---Please wait till all clients are connected: %d/%d\n",
					ConnectionCount, b.C)
				continue
			} else if !doneSuscribe {
				fmt.Print("---Already unsubscribe\n")
				continue
			}
			doneSuscribe = false
			for i := 0; i < b.C; i++ {
				go b.unsubscribe(i, topicName)
			}
			fmt.Printf("---unsubscribing to topic: %s, please wait\n", topicName)
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
		case "tst1":
			if len(ss) < 2 {
				fmt.Print("Must provide topicName, e.g.: sc yy1\n")
				continue
			}
			go b.testSubscribeThenUnsubscribe(ss[1])
		case "h":
			fallthrough
		case "?":
			fmt.Fprint(os.Stderr, fmt.Sprintf(commands))
		default:
			fmt.Println("unknow command")
		}
	}
}
