package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ewangplay/eventbus-sdk-go"
)

func main() {
	var pNodes = flag.String("nodes", "172.16.199.8:4150,172.16.199.8:5150,172.16.199.8:6150", "eventbus nodes addresses, sepatated by ','")
	var pTopic = flag.String("topic", "test01", "eventbus publish event topic")
	var pWorkerNum = flag.Int64("worker", 1000, "concurrent worker number")
	flag.Parse()

	// Set producer instance-level log options
	logger := log.New(os.Stdout, "[Producer-Test] ", log.Ldate|log.Ltime)
	sdk.SetLogger(logger, sdk.LogLevelError)

	nodeAddrs := strings.Split(*pNodes, ",")
	producers := make([]*sdk.Producer, len(nodeAddrs))

	for i, node := range nodeAddrs {
		cfg := &sdk.Config{
			Nodes: []string{node},
		}
		p, err := sdk.NewProducer(cfg)
		if err != nil {
			fmt.Printf("New producer error: %v\n", err)
			return
		}
		producers[i] = p
	}

	startTime := time.Now()

	var p *sdk.Producer
	var msg string
	var err error
	var i int64
	for i = 0; i < *pWorkerNum; i++ {
		p = producers[i%int64(len(producers))]
		msg = fmt.Sprintf("Hello, World. %d", i)

		err = p.Publish(*pTopic, []byte(msg))
		if err != nil {
			fmt.Printf("Publish message[%s: %s] error: %v\n", *pTopic, msg, err)
			return
		}
	}

	cost := time.Since(startTime)
	fmt.Printf("Cost: %v\n", cost)
	fmt.Printf("TPS: %v\n", *pWorkerNum*int64(time.Second)/int64(cost))

	for _, p := range producers {
		if p != nil {
			p.Close()
		}
	}
}
