package nsq

import (
	"fmt"
	"strings"
	"time"

	"github.com/ewangplay/eventbus/i"
	"github.com/nsqio/go-nsq"
)

const RETRY_MAX_COUNT = 3

type NSQProducer struct {
	*NSQContext
	producer *nsq.Producer
	isReady  bool
}

func NewNSQProducer(ctx *NSQContext) (*NSQProducer, error) {
	this := &NSQProducer{}

	this.NSQContext = ctx

	err := this.init()
	if err != nil {
		return nil, err
	}

	return this, nil
}

func (this *NSQProducer) Close() error {
	if this.isReady {
		if this.producer != nil {
			this.producer.Stop()
		}
		this.isReady = false
	}

	return nil
}

func (this *NSQProducer) init() error {

	nsqdNetAddr := fmt.Sprintf("127.0.0.1:%s", strings.Split(this.opts.NSQTCPAddress, ":")[1])
	cfg := nsq.NewConfig()
	p, err := nsq.NewProducer(nsqdNetAddr, cfg)
	if err != nil {
		this.Error("New nsq producer connect to %s error: %v", nsqdNetAddr, err)
		return err
	}

	this.producer = p
	this.isReady = true

	return nil
}

func (this *NSQProducer) getProducer() *nsq.Producer {
	if !this.isReady {
		if err := this.init(); err != nil {
			panic("nsq producer init fail!")
		}
	}

	return this.producer
}

func (this *NSQProducer) Publish(msg i.IMessage) (err error) {
	if !this.isReady {
		return fmt.Errorf("nsq producer instance is not ready")
	}

	retry_count := RETRY_MAX_COUNT
	for {
		err = this.getProducer().Publish(msg.GetSubject(), msg.GetData())
		if err != nil {
			this.Error("Publish message[%s:%s] to nsqd error: %v", msg.GetSubject(), msg.GetData(), err)

			if strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "connection refused") {
				this.Error("Connection exception, try again... [%d] times", RETRY_MAX_COUNT-retry_count+1)

				this.Close()

				if retry_count > 0 {
					time.Sleep(3 * time.Second)
					retry_count--
					continue
				}
			}
		}

		//succ or fail
		break
	}

	if err != nil {
		this.Error("Publish message[%s:%s] to nsqd error: %v", msg.GetSubject(), msg.GetData(), err)
	} else {
		this.Info("Publish message[%s:%s] to nsqd succ", msg.GetSubject(), msg.GetData())
	}

	return
}
