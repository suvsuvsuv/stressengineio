package boomer

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
  pm    push a message
  set   set api host, 1: xuduo; 2: liushihai
  tst1  run test 1
`

const inputdelimiter = '\n'

func (b *Boomer) readConsole() {
	doneSuscribe := false
	topicName := "yy1"
	ApiHost = "http://test.mlc.yy.com/api/bcproxy_svr/broadcast?appId=100001&sign=&data={%22topic%22:%22TOPIC_NAME%22,%22message%22:%22123%22,%22alive_time%22:12,%22is_order%22:true}"
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
				go b.unsubscribe(i, topicName, nil)
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
		case "set":
			if len(ss) < 2 {
				fmt.Print("Must provide #, e.g.: set 1\n")
				continue
			}
			i, err := strconv.Atoi(ss[1])
			if err != nil {
				fmt.Printf("# must be digit, now is %s\n", ss[1])
				continue
			}
			if i == 1 {
				ApiHost = "http://wstest.yy.com:11001/api/push?pushAll=true&topic=TOPIC_NAME&json=123&timeToLive=12"
				fmt.Print("Set to xuduo's api server\n")
			} else if i == 2 {
				ApiHost =
					"http://test.mlc.yy.com/api/bcproxy_svr/broadcast?appId=100001&sign=&data={%22topic%22:%22TOPIC_NAME%22,%22message%22:%22123%22,%22alive_time%22:12,%22is_order%22:true}"
				fmt.Print("Set to Liushihai's api server\n")
			}
		case "pm":
			pushMsgToServer(topicName)
		case "tst1":
			if len(ss) < 2 {
				fmt.Print("Must provide topicName, e.g.: tst1 yy1\n")
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
