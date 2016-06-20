package boomer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const MAX_TIMES = 5

func (b *Boomer) testSubscribeThenUnsubscribe(topicName string) {
	// sc -> send a msg -> usc
	var doneSuscribeWg sync.WaitGroup
	pushIDSent := false
	SubscribeCount = 0

	fmt.Printf("\n---Topic name: %v\n", topicName)
	startMessageCount()

	totalAvgDuration := float64(0)
	maxDuration := float64(0)
	minDuration := float64(100) // big enough value of delay
	for times := 0; times < MAX_TIMES; times++ {
		avgDuration := float64(0)
		//1. sc
		doneSuscribeWg.Add(b.C)
		for i := 0; i < b.C; i++ {
			go func(idx int) {
				if !pushIDSent {
					b.sendPushID(idx)
				}
				b.subscribe(idx, topicName, &doneSuscribeWg)
			}(i)
		}

		doneSuscribeWg.Wait()
		if !pushIDSent {
			pushIDSent = true
		}
		fmt.Printf("\n---Subscibing done: %v\n", SubscribeCount)
		//send 1 msg
		now := time.Now()
		for i := 0; i < b.C; i++ {
			b.startTimes[i] = now
		}
		pushMsgToServer(topicName)
		//duration在runWorkerEngineIo() － case "Message"里计算

		time.Sleep(10 * time.Second)
		totalValideCount := 0
		totalDuration := float64(0)
		for i := 0; i < b.C; i++ {
			if b.durations[i] <= 0 {
				fmt.Printf("\n---Invalid duration: #%d/%d\n", i, b.durations[1])
				continue
			}
			currentDuration := b.durations[i].Seconds()
			totalDuration += currentDuration
			totalValideCount++
			if maxDuration < currentDuration {
				maxDuration = currentDuration
			}
			if minDuration > currentDuration {
				minDuration = currentDuration
			}
		}
		avgDuration = totalDuration / float64(totalValideCount)
		totalAvgDuration += avgDuration
		// unsubscribe
		doneSuscribeWg.Add(b.C)
		for i := 0; i < b.C; i++ {
			go b.unsubscribe(i, topicName, &doneSuscribeWg)
		}
		doneSuscribeWg.Wait()
		fmt.Printf("---Unsubscibing done: %v\n", SubscribeCount)
		fmt.Printf("---%d: Received message count: %d---\n", times, MessageCount)
		fmt.Printf("---%d: Current Avg duration: %f---\n", times, avgDuration)
		time.Sleep(1 * time.Second)
	}
	totalAvgDuration /= float64(MAX_TIMES)
	fmt.Printf("---: Total Avg duration: %f, Max: %f, Min: %f---\n", totalAvgDuration,
		maxDuration, minDuration)
}

const SUCCESS_RESPONSE = "\"code\":0"

var ApiHost string

func pushMsgToServer(topicName string) string {
	// url: http://test.mlc.yy.com/api/bcproxy_svr/broadcast?appId=100001&sign=&data={%22topic%22:%22yy3%22,%22message%22:%22123%22,%22alive_time%22:12,%22is_order%22:true}
	client := &http.Client{}
	urlStr := strings.Replace(ApiHost, "TOPIC_NAME", topicName, 2)
	fmt.Printf("---Push message: %s\n", urlStr)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Fatalln(err)
	}

	//req.Header.Set("{0} {1} HTTP/1.1","GET")
	req.Header.Set("Connection", "close")
	//req.Header.Set("Host:","www.test.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent",
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.2; SV1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)")
	req.Header.Set("Accept-Language", "zh-CN")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("http send error: %v\n", err)
		return ""
	}
	bodyStr := string(body)
	/*if (strings.Contains(bodyStr, SUCCESS_RESPONSE) {
		return
	} */
	return bodyStr
}
