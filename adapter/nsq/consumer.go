package nsq

import (
	"fmt"

	comm "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/i"
	"github.com/nsqio/go-nsq"
)

type Consumer struct {
	*Context
	consumer  *nsq.Consumer
	isRunning bool
	messages  chan i.Message
}

func NewConsumer(ctx *Context, topic string) (*Consumer, error) {
	c := &Consumer{}

	c.Context = ctx

	c.messages = make(chan i.Message)

	go c.subscribe(topic, c.messages)

	return c, nil
}

func (c *Consumer) Close() error {
	c.Info("nsq consumer will be stopped")
	if c.isRunning {
		if c.consumer != nil {
			c.consumer.Stop()
		}
	}
	close(c.messages)
	return nil
}

func (c *Consumer) GetMessage() <-chan i.Message {
	return c.messages
}

func (c *Consumer) subscribe(topic string, messages chan<- i.Message) error {
	var outputStr string
	var err error

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = c.opts.NSQMaxInFlight
	c.consumer, err = nsq.NewConsumer(topic, topic, cfg)
	if err != nil {
		c.Error("New nsq consumer[%s:%s] error: %v", topic, topic, err)
		return err
	}
	c.consumer.AddHandler(&msgHandler{topic, messages})

	if c.opts.NSQCluster {
		err = c.consumer.ConnectToNSQLookupds(c.opts.NSQLookupdTCPAddresses)
		if err == nil {
			c.Info("Connect to nsqlookupd service[%s] succ", c.opts.NSQLookupdTCPAddresses)
			c.isRunning = true
		} else {
			c.Error("Connect to nsqlookupd service[%s] error: %v", c.opts.NSQLookupdTCPAddresses, err)
		}
	}

	if !c.isRunning {
		err = c.consumer.ConnectToNSQD(c.opts.NSQTCPAddress)
		if err == nil {
			c.Info("Connect to nsqd service[%s] succ", c.opts.NSQTCPAddress)
			c.isRunning = true
		} else {
			c.Error("Connect to nsqd service[%s] error: %v", c.opts.NSQTCPAddress, err)
		}
	}

	if c.isRunning {
		//waiting current consumer to stop
		<-c.consumer.StopChan

		c.Info("NSQ Consumer [%s/%s] exit...", topic, topic)

		c.isRunning = false

	} else {
		outputStr = fmt.Sprintf("NSQ Consumer [%s/%s] creation failed", topic, topic)
		c.Error(outputStr)

		if c.consumer != nil {
			c.consumer.Stop()
			c.consumer = nil
		}
		return fmt.Errorf(outputStr)
	}

	return nil
}

type msgHandler struct {
	topic    string
	messages chan<- i.Message
}

func (h *msgHandler) HandleMessage(m *nsq.Message) error {
	h.messages <- &comm.EBMessage{Subject: h.topic, Data: m.Body}
	return nil
}
