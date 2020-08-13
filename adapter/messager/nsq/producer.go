package nsq

import (
	"fmt"
	"strings"
	"time"

	"github.com/ewangplay/eventbus/i"
	"github.com/nsqio/go-nsq"
)

const RetryMaxCount = 3

type Producer struct {
	*Context
	producer *nsq.Producer
	isReady  bool
}

func NewProducer(ctx *Context) (*Producer, error) {
	p := &Producer{}

	p.Context = ctx

	err := p.init()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Producer) Close() error {
	if p.isReady {
		if p.producer != nil {
			p.producer.Stop()
		}
		p.isReady = false
	}

	return nil
}

func (p *Producer) init() error {

	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(p.opts.NSQTCPAddress, cfg)
	if err != nil {
		p.Error("New nsq producer connect to %s error: %v", p.opts.NSQTCPAddress, err)
		return err
	}

	p.producer = producer
	p.isReady = true

	return nil
}

func (p *Producer) getProducer() *nsq.Producer {
	if !p.isReady {
		if err := p.init(); err != nil {
			panic("nsq producer init fail!")
		}
	}

	return p.producer
}

func (p *Producer) Publish(msg i.Message) (err error) {
	if !p.isReady {
		return fmt.Errorf("nsq producer instance is not ready")
	}

	retryCount := RetryMaxCount
	for {
		err = p.getProducer().Publish(msg.GetSubject(), msg.GetData())
		if err != nil {
			p.Error("Publish message[%s:%s] to nsqd error: %v", msg.GetSubject(), msg.GetData(), err)

			if strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "connection refused") {
				p.Error("Connection exception, try again... [%d] times", RetryMaxCount-retryCount+1)

				p.Close()

				if retryCount > 0 {
					time.Sleep(3 * time.Second)
					retryCount--
					continue
				}
			}
		}

		//succ or fail
		break
	}

	if err != nil {
		p.Error("Publish message[%s:%s] to nsqd error: %v", msg.GetSubject(), msg.GetData(), err)
	} else {
		p.Info("Publish message[%s:%s] to nsqd succ", msg.GetSubject(), msg.GetData())
	}

	return
}
