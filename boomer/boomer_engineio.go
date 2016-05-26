package boomer

import (
	"fmt"
	"log"
	"strconv"
	"stressengineio/engineio2"
	"sync/atomic"
)

var MessageCount int64
var ConnectionCount int64
var ShowReceivedMessages bool
var StartCountMessages bool

// sunny :Special handling for web socket connection
func (b *Boomer) runWorkerEngineIo(n int) {
	ShowReceivedMessages = false
	StartCountMessages = false
	client, err := engineioclient2.Dial(b.RawURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	} else {
		//log.Print("connected")
		//timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		//client.SendMessage(&engineioclient.SendingObj{Type: "echo", Msg: timestamp})
		b.client = client
	}

	//Loop:
	for {
		select {
		case ev := <-client.Event:
			switch ev.Type {
			case "Open":
				atomic.AddInt64(&ConnectionCount, 1)
				if int(ConnectionCount) == b.C {
					fmt.Printf("\nConnected: %d/%d\n", ConnectionCount, b.C)
				}
			case "Close":
				atomic.AddInt64(&ConnectionCount, -1)
				fmt.Printf("Connected -1:%d\n", ConnectionCount)
				b.client = nil
			case "Message":
				if StartCountMessages {
					atomic.AddInt64(&MessageCount, 1)
				}
				if ShowReceivedMessages {
					log.Print("Got a message")
				}
			}
			//default: bug fix: do not use default in select-for, which will cause cpu-usage
		}
	}
}

func (b *Boomer) subscribe(clientIdx int, topicName string) {
	if b.client == nil {
		return
	}
	deviceID := "boom" + strconv.Itoa(clientIdx)
	v := &struct {
		ID            string                 `json:"id"`
		Version       int                    `json:"version"`
		Platform      string                 `json:"platform"`
		Topics        []string               `json:"topics"`
		LastPacketIds map[string]interface{} `json:"lastPacketIds"`
		LastUnicastID string                 `json:"lastUnicastId"`
	}{
		ID:            deviceID,
		Version:       2,
		Platform:      "stressEngineIO",
		Topics:        []string{},
		LastPacketIds: nil,
		LastUnicastID: ""}
	err := b.client.Emit("pushId", v)
	if err != nil {
		log.Fatal("emit failed:", err)
	}
	//log.Printf("new client deviceId: %s", deviceID)

	t := struct {
		Topic string `json:"topic"`
	}{topicName}
	b.client.Emit("subscribeTopic", &t)
}
