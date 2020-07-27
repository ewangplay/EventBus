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
	nsqProducer := &NSQProducer{}

	nsqProducer.NSQContext = ctx

	err := nsqProducer.init()
	if err != nil {
		return nil, err
	}

	return nsqProducer, nil
}

func (nsqProducer *NSQProducer) Close() error {
	if nsqProducer.isReady {
		if nsqProducer.producer != nil {
			nsqProducer.producer.Stop()
		}
		nsqProducer.isReady = false
	}

	return nil
}

func (nsqProducer *NSQProducer) init() error {

	cfg := nsq.NewConfig()
	p, err := nsq.NewProducer(nsqProducer.opts.NSQTCPAddress, cfg)
	if err != nil {
		nsqProducer.Error("New nsq producer connect to %s error: %v", nsqProducer.opts.NSQTCPAddress, err)
		return err
	}

	nsqProducer.producer = p
	nsqProducer.isReady = true

	return nil
}

func (nsqProducer *NSQProducer) getProducer() *nsq.Producer {
	if !nsqProducer.isReady {
		if err := nsqProducer.init(); err != nil {
			panic("nsq producer init fail!")
		}
	}

	return nsqProducer.producer
}

func (nsqProducer *NSQProducer) Publish(msg i.IMessage) (err error) {
	if !nsqProducer.isReady {
		return fmt.Errorf("nsq producer instance is not ready")
	}

	retry_count := RETRY_MAX_COUNT
	for {
		err = nsqProducer.getProducer().Publish(msg.GetSubject(), msg.GetData())
		if err != nil {
			nsqProducer.Error("Publish message[%s:%s] to nsqd error: %v", msg.GetSubject(), msg.GetData(), err)

			if strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "connection refused") {
				nsqProducer.Error("Connection exception, try again... [%d] times", RETRY_MAX_COUNT-retry_count+1)

				nsqProducer.Close()

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
		nsqProducer.Error("Publish message[%s:%s] to nsqd error: %v", msg.GetSubject(), msg.GetData(), err)
	} else {
		nsqProducer.Info("Publish message[%s:%s] to nsqd succ", msg.GetSubject(), msg.GetData())
	}

	return
}
