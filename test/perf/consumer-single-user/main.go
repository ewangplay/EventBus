package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ewangplay/eventbus-sdk-go"
)

func main() {
	var pNodes = flag.String("nodes", "172.16.199.8:4150,172.16.199.8:5150,172.16.199.8:6150", "eventbus nodes addresses, sepatated by ','")
	var pTopic = flag.String("topic", "test01", "eventbus publish event topic")
	//var pWorkerNum = flag.Int64("worker", 1000, "concurrent worker number")
	flag.Parse()

	// Set sdk-level log options
	logger := log.New(os.Stdout, "[Consumer-Node-Test] ", log.Ldate|log.Ltime)
	sdk.SetLogger(logger, sdk.LogLevelError)

	nodeAddrs := strings.Split(*pNodes, ",")
	consumers := make([]*sdk.Consumer, len(nodeAddrs))

	for i, node := range nodeAddrs {
		cfg := &sdk.Config{
			Nodes: []string{node},
		}

		c, err := sdk.NewConsumer(*pTopic, *pTopic, cfg)
		if err != nil {
			fmt.Printf("New consumer error: %v\n", err)
			return
		}
		consumers[i] = c
	}

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for _, c := range consumers {
		ch, err := c.Consume()
		if err != nil {
			fmt.Printf("Consume from eventbus error: %v\n", err)
			return
		}

		go func(ch <-chan *sdk.Message) {
			for msg := range ch {
				fmt.Printf("Receive message: %s\n", msg.Body)
			}
		}(ch)
	}

	<-exitChan

	fmt.Println("exit...")

	for _, c := range consumers {
		if c != nil {
			c.Close()
		}
	}
}
