package boomer

import (
	"fmt"
	"log"
	"strconv"
	"stressengineio/engineio2"
	"sync/atomic"
)

// pointer of engineio client
var MessageCount int64
var ConnectionCount int64
var ShowReceivedMessages bool
var StartCountMessages bool
var SubscribeCount int64

// sunny :Special handling for web socket connection
func (b *Boomer) runWorkerEngineIo(n int) {

	ShowReceivedMessages = false
	StartCountMessages = false
	client, err := engineioclient2.Dial(b.RawURL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	//Loop:
	for {
		select {
		case ev := <-client.Event:
			switch ev.Type {
			case "Open":
				b.clients[n] = client
				atomic.AddInt64(&ConnectionCount, 1)
				if int(ConnectionCount) == b.C {
					fmt.Printf("\n---Connected: %d/%d\n", ConnectionCount, b.C)
				}
			case "Close":
				atomic.AddInt64(&ConnectionCount, -1)
				log.Printf("---Disconnected by remote server:%d\n", ConnectionCount)
				b.clients[n] = nil
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
	if b.clients[clientIdx] == nil {
		log.Fatal("connnection failed:", clientIdx)
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
	err := b.clients[clientIdx].Emit("pushId", v)
	if err != nil {
		log.Fatal("emit failed:", err)
	}
	//log.Printf("new client deviceId: %s", deviceID)

	t := struct {
		Topic string `json:"topic"`
	}{topicName}
	b.clients[clientIdx].Emit("subscribeTopic", &t)
	atomic.AddInt64(&SubscribeCount, 1)
	if int(SubscribeCount) == b.C {
		fmt.Printf("\n---Subscibing done: %v\n", SubscribeCount)
	}
}
