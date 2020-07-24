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
	var pSeekerAddrs = flag.String("seeker-address", "127.0.0.1:4161", "eventbus seeker address list, separated by ','")
	var pTopic = flag.String("topic", "", "topic to be consumed")
	var pChannel = flag.String("channel", "", "channel for topic, can be empty")
	var pMaxInFlight = flag.Int("max-in-flight", 1000, "max in flight for consumer")
	flag.Parse()

	//Topic must not be empty
	if *pTopic == "" {
		flag.Usage()
		return
	}

	seekers := strings.Split(*pSeekerAddrs, ",")

	// Set sdk-level log options
	logger := log.New(os.Stdout, "[Consumer-Seeker-Test] ", log.Ldate|log.Ltime)
	sdk.SetLogger(logger, sdk.LogLevelInfo)

	// Create new Consumer
	cfg := &sdk.Config{
		Seekers:     seekers,
		MaxInFlight: *pMaxInFlight,
	}
	c, err := sdk.NewConsumer(*pTopic, *pChannel, cfg)
	if err != nil {
		fmt.Printf("New consumer error: %v\n", err)
		return
	}

	// Consume messages from eventbus
	ch, err := c.Consume()
	if err != nil {
		fmt.Printf("Consume from eventbus error: %v\n", err)
		return
	}

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func(ch <-chan *sdk.Message) {
		total := 0
		for _ = range ch {
			total += 1
			//fmt.Printf("Receive message: %s\n", msg.Body)
			fmt.Printf("topic: %v, total: %d\n", *pTopic, total)
		}
	}(ch)

	<-exitChan

	fmt.Println("exit...")

	c.Close()
}
