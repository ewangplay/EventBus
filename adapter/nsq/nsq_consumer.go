package nsq

import (
	"fmt"

	c "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/i"
	"github.com/nsqio/go-nsq"
)

type NSQConsumer struct {
	*NSQContext
	consumer  *nsq.Consumer
	isRunning bool
	messages  chan i.IMessage
}

func NewNSQConsumer(ctx *NSQContext, topic string) (*NSQConsumer, error) {
	this := &NSQConsumer{}

	this.NSQContext = ctx

	this.messages = make(chan i.IMessage)

	go this.subscribe(topic, this.messages)

	return this, nil
}

func (this *NSQConsumer) Close() error {
	this.Info("nsq consumer will be stopped")
	if this.isRunning {
		if this.consumer != nil {
			this.consumer.Stop()
		}
	}
	close(this.messages)
	return nil
}

func (this *NSQConsumer) GetMessage() <-chan i.IMessage {
	return this.messages
}

func (this *NSQConsumer) subscribe(topic string, messages chan<- i.IMessage) error {
	var outputStr string
	var err error

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = this.opts.NSQMaxInFlight
	this.consumer, err = nsq.NewConsumer(topic, topic, cfg)
	if err != nil {
		this.Error("New nsq consumer[%s:%s] error: %v", topic, topic, err)
		return err
	}
	this.consumer.AddHandler(&msgHandler{topic, messages})

	if this.opts.NSQCluster {
		err = this.consumer.ConnectToNSQLookupds(this.opts.NSQLookupdTCPAddresses)
		if err == nil {
			this.Info("Connect to nsqlookupd service[%s] succ", this.opts.NSQLookupdTCPAddresses)
			this.isRunning = true
		} else {
			this.Error("Connect to nsqlookupd service[%s] error: %v", this.opts.NSQLookupdTCPAddresses, err)
		}
	}

	if !this.isRunning {
		err = this.consumer.ConnectToNSQD(this.opts.NSQTCPAddress)
		if err == nil {
			this.Info("Connect to nsqd service[%s] succ", this.opts.NSQTCPAddress)
			this.isRunning = true
		} else {
			this.Error("Connect to nsqd service[%s] error: %v", this.opts.NSQTCPAddress, err)
		}
	}

	if this.isRunning {
		//waiting current consumer to stop
		<-this.consumer.StopChan

		this.Info("NSQ Consumer [%s/%s] exit...", topic, topic)

		this.isRunning = false

	} else {
		outputStr = fmt.Sprintf("NSQ Consumer [%s/%s] creation failed", topic, topic)
		this.Error(outputStr)

		if this.consumer != nil {
			this.consumer.Stop()
			this.consumer = nil
		}
		return fmt.Errorf(outputStr)
	}

	return nil
}

type msgHandler struct {
	topic    string
	messages chan<- i.IMessage
}

func (h *msgHandler) HandleMessage(m *nsq.Message) error {
	h.messages <- &c.EB_Message{Subject: h.topic, Data: m.Body}
	return nil
}
