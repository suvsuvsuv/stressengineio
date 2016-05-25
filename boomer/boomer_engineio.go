package boomer

import (
	"fmt"
	"log"
	"stressengineio/engineio2"
	"sync/atomic"
)

var MessageCount int64
var ConnectionCount int64

// sunny :Special handling for web socket connection
func (b *Boomer) runWorkerEngineIo(n int) {
	//var opened = false
	/*var throttle <-chan time.Time
	if b.Qps > 0 {
		throttle = time.Tick(time.Duration(1e6/(b.Qps)) * time.Microsecond)
	} */
	// b.RawURL

	/*opt := make(map[string]string)
	opt["transport"] = "websocket"
	client, err := engineioclient.Dial(
		b.RawURL,
		//"http://221.228.83.182:18301",
		opt)
	if err != nil {
		log.Fatal("dial:", err)
		//fmt.Print("error")
	} else {
		//log.Print("connected")
		//timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		//client.SendMessage(&engineioclient.SendingObj{Type: "echo", Msg: timestamp})
	}

	//Loop:
	for {
		select {
		case ev := <-client.Event:
			if !opened {
				handleNotOpen(ev, client, &opened)
			} else {
				handleOpened(ev, client)
				//break Loop
			}
		default:
		}
	} */
	client, err := engineioclient2.Dial(b.RawURL, nil)
	if err != nil {
		log.Panic("io connect %s", err)
	}
	if err != nil {
		log.Fatal("dial:", err)
		//fmt.Print("error")
	} else {
		//opened = true
		//log.Print("connected")
		//timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		//client.SendMessage(&engineioclient.SendingObj{Type: "echo", Msg: timestamp})
	}

	//Loop:
	for {
		select {
		case ev := <-client.Event:
			switch ev.Type {
			case "Open":
				atomic.AddInt64(&ConnectionCount, 1)
				if int(ConnectionCount) == b.C {
					fmt.Printf("Connected: %d/%d\n", ConnectionCount, b.C)
				}
			case "Close":
				atomic.AddInt64(&ConnectionCount, -1)
				fmt.Printf("Connected -1:%d\n", ConnectionCount)
			case "Message":
				//fmt.Print("got a message")
				//atomic.AddInt64(&MessageCount, 1)
			}
			//default: bug fix: do not use default in select-for, which will cause cpu-usage
			//fmt.Printf("Connected :%d\n", connectionCount)
		}
	}
}

/*func handleNotOpen(ev *engineioclient.Event, client *engineioclient.Client, opened *bool) {
	if !*opened && ev.Type != "Open" {
		//test.Errorf("The first event should be open!")
		log.Fatal("The first event should be open!")
	}
	*opened = true
	atomic.AddInt64(&connectionCount, 1)
	fmt.Printf("Connection :%d\n", connectionCount)
}

func handleOpened(ev *engineioclient.Event, client *engineioclient.Client) {
	//if ev.Type != "pong" {
	//	log.Fatal("The next event should be pong!")
	//}
	switch ev.Type {
	case "Message":
		//fmt.Print("got a message")
		atomic.AddInt64(&messageCount, 1)
	}
} */
