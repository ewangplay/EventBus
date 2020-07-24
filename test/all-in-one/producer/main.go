package main

import (
	"flag"
	"fmt"
	"github.com/ewangplay/eventbus-sdk-go"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var pNodeAddr = flag.String("node-address", "127.0.0.1:4150", "eventbus node address")
	var pTopics = flag.String("topic", "", "topic to be published")
	var pProducerNum = flag.Int("producer-count", 1, "count of producers, can be empty")
	var pMessagecount = flag.Int("message-count", 100, "number of messages that each topic is sent on")
	flag.Parse()

	//Topic must not be empty
	if *pTopics == "" {
		flag.Usage()
		return
	}

	topics := strings.Split(*pTopics, ",")

	cfg := &sdk.Config{
		Nodes: []string{*pNodeAddr},
	}

	for i := 0; i < *pProducerNum; i++ {
		// Create new Producer instance
		wg.Add(1)
		go func() {
			defer wg.Done()
			p, err := sdk.NewProducer(cfg)
			if err != nil {
				fmt.Printf("New producer error: %v\n", err)
				return
			}
			defer p.Close()
			fmt.Printf("Producer create succeed.")

			for _, topic := range topics {
				for pData := 1; pData <= *pMessagecount; pData++ {
					pDataStr := fmt.Sprintf("testdata%d", pData)
					err = p.Publish(topic, []byte(pDataStr))
					if err != nil {
						fmt.Printf("Publish message[%s: %s] error: %v\n", topic, pDataStr, err)
						continue
					}

					fmt.Printf("Publish message[%s: %s] succ\n", topic, pDataStr)
				}

			}
		}()
		wg.Wait()
	}
}
