package boomer

import (
	"fmt"
	"log"
	"strconv"
	"stressengineio/engineio2"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"math/rand"
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
				b.sendPushID(n)
			case "Close":
				atomic.AddInt64(&ConnectionCount, -1)
				log.Printf("---Disconnected by remote server:%d\n", ConnectionCount)
				b.clients[n] = nil
			case "Message":
				//忽略push id返回消息
				evDataStr := string(ev.Data)
				if strings.Contains(evDataStr, "pushId") {
					continue
				}
				if StartCountMessages {
					atomic.AddInt64(&MessageCount, 1)
					b.durations[n] = time.Now().Sub(b.startTimes[n])
				}
				if ShowReceivedMessages {
					log.Printf("Got a Message %s", string(ev.Data))
					//log.Print("Got a message")
				}
			}
			//default: bug fix: do not use default in select-for, which will cause cpu-usage very high
		}
	}
}

func (b *Boomer) sendPushID(clientIdx int) {
	if b.clients[clientIdx] == nil {
		log.Fatal("connnection failed:", clientIdx)
		return
	}

	idPrefix := b.PushIdPrefix
	if idPrefix == ""  || len(idPrefix) < 16 {
         idPrefix = string(Krand(16, KC_RAND_KIND_ALL))
	}

	deviceID := idPrefix + strconv.Itoa(clientIdx)

	v := &struct {
		ID            string                 `json:"id"`
		Version       int                    `json:"version"`
		Platform      string                 `json:"platform"`
		Topics        []string               `json:"topics"`
		LastPacketIds map[string]interface{} `json:"lastPacketIds"`
		LastUnicastID string                 `json:"lastUnicastId"`
	}{
		ID:            deviceID,
		Version:       1,
		Platform:      "stressEngineIO",
		Topics:        []string{},
		LastPacketIds: nil,
		LastUnicastID: ""}
	err := b.clients[clientIdx].Emit("pushId", v)
	if err != nil {
		log.Fatal("emit failed:", err)
	}
	//log.Printf("new client deviceId: %s", deviceID)
}

func (b *Boomer) addTag(clientIdx int, tagName string) {
	if b.clients[clientIdx] == nil {
		log.Fatal("connnection failed:", clientIdx)
		return
	}

	t := struct {
		Topic string `json:"tag"`
	}{tagName}
	b.clients[clientIdx].Emit("addTag", &t)
}

func (b *Boomer) removeTag(clientIdx int, tagName string) {
	if b.clients[clientIdx] == nil {
		log.Fatal("connnection failed:", clientIdx)
		return
	}

	t := struct {
		Topic string `json:"tag"`
	}{tagName}
	b.clients[clientIdx].Emit("removeTag", &t)
}

func (b *Boomer) subscribe(clientIdx int, topicName string, doneSuscribeWg *sync.WaitGroup) {
	defer func() {
		if doneSuscribeWg != nil {
			(*doneSuscribeWg).Done()
		}
	}()
	t := struct {
		Topic string `json:"topic"`
	}{topicName}
	b.clients[clientIdx].Emit("subscribeTopic", &t)
	atomic.AddInt64(&SubscribeCount, 1)
	if int(SubscribeCount) == b.C {
		if doneSuscribeWg == nil {
			fmt.Printf("\n---Subscibing done: %v\n", SubscribeCount)
		}
	}

}

func (b *Boomer) unsubscribe(clientIdx int, topicName string, doneSuscribeWg *sync.WaitGroup) {
	defer func() {
		if doneSuscribeWg != nil {
			(*doneSuscribeWg).Done()
		}
	}()
	t := struct {
		Topic string `json:"topic"`
	}{topicName}
	b.clients[clientIdx].Emit("unsubscribeTopic", &t)
	atomic.AddInt64(&SubscribeCount, -1)
	if int(SubscribeCount) == 0 {
		if doneSuscribeWg == nil {
			fmt.Printf("\n---unsubscibing done: %v\n", SubscribeCount)
		}
	}
}

const (
    KC_RAND_KIND_NUM   = 0  // 纯数字
    KC_RAND_KIND_LOWER = 1  // 小写字母
    KC_RAND_KIND_UPPER = 2  // 大写字母
    KC_RAND_KIND_ALL   = 3  // 数字、大小写字母
)
 
// 随机字符串
func Krand(size int, kind int) []byte {
    ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
    is_all := kind > 2 || kind < 0
    rand.Seed(time.Now().UnixNano())
    for i :=0; i < size; i++ {
        if is_all { // random ikind
            ikind = rand.Intn(3)
        }
        scope, base := kinds[ikind][0], kinds[ikind][1]
        result[i] = uint8(base+rand.Intn(scope))
    }
    return result
}
