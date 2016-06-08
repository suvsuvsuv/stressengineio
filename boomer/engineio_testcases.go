package boomer

import (
	"fmt"
	"sync"
)

func (b *Boomer) testSubscribeThenUnsubscribe(topicName string) {
	// sc -> send a msg -> usc
	//1. sc
	var doneSuscribeWg sync.WaitGroup
	doneSuscribeWg.Add(b.C)
	pushIDSent := false
	SubscribeCount = 0
	for i := 0; i < b.C; i++ {
		go func(idx int) {
			if !pushIDSent {
				b.sendPushID(idx)
			}
			b.subscribe(idx, topicName, &doneSuscribeWg)
		}(i)
	}
	if !pushIDSent {
		pushIDSent = true
	}
	doneSuscribeWg.Wait()

	fmt.Printf("\n---Subscibing done 2: %v\n", SubscribeCount)
	//send 1 msg
}
