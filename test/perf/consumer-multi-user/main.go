package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ewangplay/eventbus-sdk-go"
)

var (
	gTotalMsgNum int64
)

func main() {
	var pNodes = flag.String("nodes", "172.16.199.8:4150,172.16.199.8:5150,172.16.199.8:6150", "eventbus nodes addresses, sepatated by ','")
	var pTopic = flag.String("topic", "test01", "eventbus publish event topic")
	var pUserNum = flag.Int64("user", 10, "concurrent worker number")
	var pWorkerNum = flag.Int64("worker", 10000, "concurrent worker number")
	flag.Parse()

	// Set sdk-level log options
	logger := log.New(os.Stdout, "[Consumer-Node-Test] ", log.Ldate|log.Ltime)
	sdk.SetLogger(logger, sdk.LogLevelError)

	nodeAddrs := strings.Split(*pNodes, ",")
	consumers := make([][]*sdk.Consumer, *pUserNum)

	for i := range consumers {
		consumers[i] = make([]*sdk.Consumer, len(nodeAddrs))
		for j, node := range nodeAddrs {
			cfg := &sdk.Config{
				Nodes: []string{node},
			}

			c, err := sdk.NewConsumer(*pTopic, *pTopic, cfg)
			if err != nil {
				fmt.Printf("New consumer error: %v\n", err)
				return
			}
			consumers[i][j] = c
		}
	}

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	startTime := time.Now()

	expectTotalMsgNum := *pUserNum * *pWorkerNum

	for _, cs := range consumers {
		for _, c := range cs {
			ch, err := c.Consume()
			if err != nil {
				fmt.Printf("Consume from eventbus error: %v\n", err)
				return
			}

			go func(ch <-chan *sdk.Message) {
				for msg := range ch {
					fmt.Printf("Receive message: %s\n", msg.Body)
					atomic.AddInt64(&gTotalMsgNum, 1)
					if gTotalMsgNum == expectTotalMsgNum {
						exitChan <- 1
					}
				}
			}(ch)
		}
	}

	<-exitChan

	cost := time.Since(startTime)
	fmt.Printf("Cost: %v\n", cost)
	fmt.Printf("TPS: %v\n", gTotalMsgNum*int64(time.Second)/int64(cost))

	fmt.Println("exit...")

	for _, cs := range consumers {
		for _, c := range cs {
			if c != nil {
				c.Close()
			}
		}
	}
}
