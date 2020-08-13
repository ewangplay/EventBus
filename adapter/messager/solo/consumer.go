package solo

import (
	"fmt"

	comm "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/i"
)

// Consumer represents Solo Mode Messager Consumer
type Consumer struct {
	*Context
	messages chan i.Message
}

// NewConsumer creates a new Solo Mode Messager Consumer instance
func NewConsumer(ctx *Context, topic string) (*Consumer, error) {
	c := &Consumer{}

	c.Context = ctx

	c.messages = make(chan i.Message)

	go c.subscribe(topic, c.messages)

	return c, nil
}

// Close closes the current Consumer instance
func (c *Consumer) Close() error {
	close(c.messages)
	return nil
}

// GetMessage returns the Message channel
func (c *Consumer) GetMessage() <-chan i.Message {
	return c.messages
}

func (c *Consumer) subscribe(topic string, messages chan<- i.Message) (err error) {
	c.Lock()
	queue, found := c.Queues[topic]
	if !found {
		queue = make(chan []byte, QueueMaxSize)
		c.Queues[topic] = queue
	}
	c.Unlock()

	// 用于发生panic时的恢复
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// 如果messages被关闭，继续写入数据会导致panic
	for data := range queue {
		c.messages <- &comm.EBMessage{Subject: topic, Data: data}
	}

	return nil
}
