package adapter

import (
	"fmt"

	"github.com/ewangplay/eventbus/adapter/nsq"
	"github.com/ewangplay/eventbus/adapter/rabbitmq"
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

type Messager struct {
	i.ILogger
	context   i.IContext
	producer  i.IProducer
	opts      *config.EB_Options
	consumers map[string]i.IConsumer
}

func NewMessager(opts *config.EB_Options, logger i.ILogger) (*Messager, error) {
	this := &Messager{}

	this.opts = opts
	this.ILogger = logger

	var err error
	var context i.IContext
	var producer i.IProducer
	if this.opts.NSQEnable {
		context, err = nsq.GetNSQContext(opts, logger)
		if err != nil {
			this.Error("Get nsq context instance error: %v", err)
			return nil, err
		}

		producer, err = context.CreateProducer()
		if err != nil {
			this.Error("Create nsq producer error: %v", err)
			return nil, err
		}

	} else {
		context, err = rabbitmq.GetRabbitMQContext(opts, logger)
		if err != nil {
			this.Error("Get rabbitmq context instance error: %v", err)
			return nil, err
		}

		producer, err = context.CreateProducer()
		if err != nil {
			this.Error("Create rabbitmq producer error: %v", err)
			return nil, err
		}

	}

	this.context = context
	this.producer = producer
	this.consumers = make(map[string]i.IConsumer, 1)

	return this, nil
}

func (this *Messager) Close() error {
	this.Info("Messager will close")
	this.producer.Close()
	for _, consumer := range this.consumers {
		consumer.Close()
	}
	return nil
}

func (this *Messager) Publish(msg i.IMessage) error {
	if this.producer == nil {
		return fmt.Errorf("producer not valid")
	}
	return this.producer.Publish(msg)
}

func (this *Messager) Subscribe(subject string) (<-chan i.IMessage, error) {

	//If consumer has already existed, get message chan directly
	old_consumer, ok := this.consumers[subject]
	if ok {
		return old_consumer.GetMessage(), nil
	}

	var err error
	var consumer i.IConsumer
	if this.opts.NSQEnable {

		if this.context == nil {
			this.context, err = nsq.GetNSQContext(this.opts, this.ILogger)
			if err != nil {
				this.Error("Get nsq context instance error: %v", err)
				return nil, err
			}
		}

		consumer, err = this.context.CreateConsumer(subject)
		if err != nil {
			this.Error("Create nsq consumer error: %v", err)
			return nil, err
		}

	} else {

		if this.context == nil {
			this.context, err = rabbitmq.GetRabbitMQContext(this.opts, this.ILogger)
			if err != nil {
				this.Error("Get rabbitmq context instance error: %v", err)
				return nil, err
			}
		}

		consumer, err = this.context.CreateConsumer(subject)
		if err != nil {
			this.Error("Create rabbitmq consumer error: %v", err)
			return nil, err
		}

	}

	this.consumers[subject] = consumer

	return consumer.GetMessage(), nil
}

func (this *Messager) Unsubscribe(subject string) error {
	consumer, ok := this.consumers[subject]
	if !ok {
		this.Info("Unsubscribe a consumer which does not exist: %v", subject)
		return nil
	}

	err := consumer.Close()
	if err != nil {
		this.Error("Close consumer[%v] error: %v", subject, err)
		return err
	}

	delete(this.consumers, subject)

	return nil
}
