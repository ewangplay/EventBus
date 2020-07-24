package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ewangplay/eventbus-sdk-go"
)

func main() {
	var pNodes = flag.String("nodes", "172.16.199.8:4150,172.16.199.8:5150,172.16.199.8:6150", "eventbus nodes addresses, sepatated by ','")
	var pTopic = flag.String("topic", "test01", "eventbus publish event topic")
	var pUserNum = flag.Int64("user", 10, "concurrent user number")
	var pWorkerNum = flag.Int64("worker", 10000, "concurrent worker number")
	flag.Parse()

	// Set producer instance-level log options
	logger := log.New(os.Stdout, "[Producer-Test] ", log.Ldate|log.Ltime)
	sdk.SetLogger(logger, sdk.LogLevelError)

	nodeAddrs := strings.Split(*pNodes, ",")
	producers := make([][]*sdk.Producer, *pUserNum)

	for i := range producers {
		producers[i] = make([]*sdk.Producer, len(nodeAddrs))
		for j, node := range nodeAddrs {
			cfg := &sdk.Config{
				Nodes: []string{node},
			}
			p, err := sdk.NewProducer(cfg)
			if err != nil {
				fmt.Printf("New producer error: %v\n", err)
				return
			}
			producers[i][j] = p
		}
	}

	var wg sync.WaitGroup
	var i, j int64

	startTime := time.Now()

	for i = 0; i < *pUserNum; i++ {
		for j = 0; j < *pWorkerNum; j++ {
			wg.Add(1)

			go func(p *sdk.Producer, topic string, msg string) {
				err := p.Publish(topic, []byte(msg))
				if err != nil {
					fmt.Printf("Publish message[%s: %s] error: %v\n", topic, msg, err)
				}
				wg.Done()
			}(producers[i][j%int64(len(producers[i]))], *pTopic, fmt.Sprintf("Hello, World. %d - %d", i, j))
		}
	}

	wg.Wait()

	cost := time.Since(startTime)
	fmt.Printf("Cost: %v\n", cost)
	fmt.Printf("TPS: %v\n", (*pUserNum)*(*pWorkerNum)*int64(time.Second)/int64(cost))

	for _, ps := range producers {
		for _, p := range ps {
			if ps != nil {
				p.Close()
			}
		}
	}
}
